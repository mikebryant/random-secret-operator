#!/usr/bin/env python

'''
Copyright 2018 Mike Bryant. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
'''

from http.server import BaseHTTPRequestHandler, HTTPServer
import base64
import json
import copy
import re
import secrets


def is_acceptable(
    data,
    desired_bytes,
):
    return len(data)/2 >= desired_bytes

class Controller(BaseHTTPRequestHandler):
    def sync(self, randomsecret, children):

        complete = False

        desired_field_name = 'random'
        desired_bytes = 64

        # Find existing secret, if it exists
        desired_name = randomsecret['metadata']['name']
        existing_random_data = base64.standard_b64decode(
            children['Secret.v1'].get(
                desired_name,
                {},
            ).get(
                'data',
                {},
            ).get(
                desired_field_name,
                '',
            ),
        )

        if is_acceptable(
            data=existing_random_data,
            desired_bytes=desired_bytes,
        ):
            complete = True
        else:
            existing_random_data = secrets.token_hex(
                desired_bytes,
            ).encode('utf-8')


        desired_secret = {
            'apiVersion': 'v1',
            'kind': 'Secret',
            'metadata': {
                'name': desired_name,
            },
            'data': {
                desired_field_name: base64.standard_b64encode(
                    existing_random_data,
                ).decode('utf-8'),
            },
        }

        return {
            'children': [
                desired_secret,
            ],
            'status': {
                'conditions': [
                    {
                        'status': str(complete),
                        'type': 'Complete',
                    },
                ],
            },
        }


    def do_POST(self):
        observed = json.loads(self.rfile.read(int(self.headers.get('content-length'))))
        desired = self.sync(observed['parent'], observed['children'])

        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(desired).encode('utf-8'))

HTTPServer(('', 8080), Controller).serve_forever()
