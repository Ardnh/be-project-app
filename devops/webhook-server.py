#!/usr/bin/env python3
# /opt/deploy/webhook-server.py

from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import subprocess
import os
import logging
import hmac
import hashlib
from datetime import datetime

# Configuration
WEBHOOK_SECRET = os.environ.get('WEBHOOK_SECRET')
DEPLOY_SCRIPT = '/opt/deploy/deploy.sh'
PORT = 9000

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/opt/logs/webhook.log'),
        logging.StreamHandler()
    ]
)

class DeploymentQueue:
    """Simple deployment queue to prevent concurrent deployments"""
    def __init__(self):
        self.is_deploying = False
        self.queue = []

    def add(self, payload):
        self.queue.append(payload)
        return len(self.queue)

    def start_deployment(self):
        if self.is_deploying or not self.queue:
            return None
        self.is_deploying = True
        return self.queue.pop(0)

    def finish_deployment(self):
        self.is_deploying = False

queue = DeploymentQueue()

class WebhookHandler(BaseHTTPRequestHandler):

    def log_message(self, format, *args):
        logging.info(format % args)

    def verify_signature(self, payload, signature):
        """Verify webhook signature"""
        if not signature:
            return False

        expected = hmac.new(
            WEBHOOK_SECRET.encode(),
            payload,
            hashlib.sha256
        ).hexdigest()

        return hmac.compare_digest(signature, expected)

    def do_POST(self):
        # Verify token
        token = self.headers.get('X-Webhook-Token')
        if token != WEBHOOK_SECRET:
            logging.warning(f"Invalid token from {self.client_address[0]}")
            self.send_error(403, "Invalid token")
            return

        # Read payload
        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length)

        try:
            payload = json.loads(body.decode('utf-8'))

            # Validate required fields
            required = ['repository', 'branch', 'commit', 'image']
            if not all(k in payload for k in required):
                self.send_error(400, "Missing required fields")
                return

            logging.info(f"✅ Valid webhook received")
            logging.info(f"Repository: {payload['repository']}")
            logging.info(f"Branch: {payload['branch']}")
            logging.info(f"Commit: {payload['commit'][:8]}")
            logging.info(f"Image: {payload['image']}")
            logging.info(f"Deployer: {payload.get('deployer', 'unknown')}")

            # Add to queue
            position = queue.add(payload)

            if position == 1 and not queue.is_deploying:
                # Start deployment immediately
                self._trigger_deployment(payload)
            else:
                logging.info(f"📋 Deployment queued (position: {position})")

            # Send response
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                'status': 'success',
                'message': 'Deployment triggered' if position == 1 else f'Deployment queued (position: {position})',
                'timestamp': datetime.now().isoformat(),
                'queue_position': position
            }
            self.wfile.write(json.dumps(response).encode())

        except json.JSONDecodeError:
            logging.error("Invalid JSON payload")
            self.send_error(400, "Invalid JSON")
        except Exception as e:
            logging.error(f"Error: {str(e)}")
            self.send_error(500, str(e))

    def _trigger_deployment(self, payload):
        """Trigger deployment in background"""
        try:
            # Pass payload as environment variables
            env = os.environ.copy()
            env.update({
                'DEPLOY_IMAGE': payload['image'],
                'DEPLOY_BRANCH': payload['branch'],
                'DEPLOY_COMMIT': payload['commit'],
                'DEPLOY_REPOSITORY': payload['repository'],
                'DEPLOYER': payload.get('deployer', 'unknown')
            })

            # Start deployment process
            subprocess.Popen(
                ['/bin/bash', DEPLOY_SCRIPT],
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                start_new_session=True
            )

            logging.info("🚀 Deployment process started")

        except Exception as e:
            logging.error(f"Failed to trigger deployment: {str(e)}")
            queue.finish_deployment()

    def do_GET(self):
        """Health check & status endpoint"""
        if self.path == '/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                'status': 'healthy',
                'is_deploying': queue.is_deploying,
                'queue_length': len(queue.queue)
            }
            self.wfile.write(json.dumps(response).encode())
        elif self.path == '/status':
            # Read recent logs
            try:
                with open('/opt/logs/deploy.log', 'r') as f:
                    recent_logs = f.readlines()[-50:]

                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                response = {
                    'is_deploying': queue.is_deploying,
                    'queue_length': len(queue.queue),
                    'recent_logs': recent_logs
                }
                self.wfile.write(json.dumps(response).encode())
            except Exception as e:
                self.send_error(500, str(e))
        else:
            self.send_error(404)

if __name__ == '__main__':
    if not WEBHOOK_SECRET:
        logging.error("WEBHOOK_SECRET not set!")
        exit(1)

    server = HTTPServer(('0.0.0.0', PORT), WebhookHandler)
    logging.info(f'🎣 Webhook server started on port {PORT}')

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logging.info('Shutting down...')
        server.shutdown()
