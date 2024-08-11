import threading
from http.server import HTTPServer, BaseHTTPRequestHandler
import json

class HealthHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            # Use a dictionary with unique keys for each entry.
            health_status = {'status': 'healthy', 'message': 'hello'}
            self.wfile.write(json.dumps(health_status).encode('utf-8'))
        else:
            self.send_response(404)
            self.end_headers()

def run_server(port):
    server = HTTPServer(('localhost', port), HealthHandler)
    print(f'Server running on port {port}')
    server.serve_forever()

if __name__ == '__main__':
    ports = [8080, 8081, 8082, 8083]
    threads = []

    for port in ports:
        thread = threading.Thread(target=run_server, args=(port,))
        thread.start()
        threads.append(thread)

    for thread in threads:
        thread.join()