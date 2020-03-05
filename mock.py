# -*- coding: utf-8 -*-
#
# Copyright © 2020 white <white@Whites-Mac-Air.local>
#
# Distributed under terms of the MIT license.

"""
"""

from flask import Flask, request, jsonify


app = Flask('xpra-cmd')


@app.route('/s3/user/server/users', methods=['POST'])
def launch():
    print(request.headers)
    data = request.get_json()
    print(data)
    rt = {
            "status": 0,
            "msg": "success"
            }
    return jsonify(rt)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000, debug=True)

