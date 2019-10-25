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
import logging
import re
import secrets

logging.basicConfig(level=logging.DEBUG)
LOGGER = logging.getLogger(__name__)

def is_acceptable(
    data,
    desired_bytes,
    metadata,
):
    if len(data) == desired_bytes:
        return True
    else:
        LOGGER.info(
            '%s/%s: Existing length %d != desired length %d',
            metadata['namespace'],
            metadata['name'],
            len(data),
            desired_bytes,
        )
        return False

class Controller(BaseHTTPRequestHandler):
    def sync(self, randomsecret, children):

        LOGGER.debug(
            'Reconciling %s/%s...',
            randomsecret['metadata']['namespace'],
            randomsecret['metadata']['name'],
        )

        complete = False

        spec = randomsecret.get('spec', {})

        desired_field_name = 'random'
        desired_length = spec.get('length', 128)

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
            desired_bytes=desired_length,
            metadata=randomsecret['metadata'],
        ):
            complete = True
        else:
            existing_random_data = secrets.token_hex(
                desired_length,
            ).encode('utf-8')[:desired_length]


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
