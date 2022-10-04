import os
from flask import Flask
from flask import send_from_directory
from flask import request
from flask import make_response
import socket
from waitress import serve
import hashlib
import json
import logging
import datetime
import random
import requests
import uuid

logging.basicConfig(
        level = logging.DEBUG,
        format='[%(asctime)s][%(filename)s][%(levelname)s][%(message)s]',
        datefmt='%Y-%m-%d %H:%M:%S',
        )

app = Flask(__name__)
ip = None
port = None
id = None


@app.route("/status", methods=['GET'])
def status():
    print("被请求一次")
    return "ok", 200

def Init():
    global ip
    global port
    global id
    
    hostname = socket.gethostname()
    ip = socket.gethostbyname(hostname) # 获取本地Ip
    port = 9900 + random.randint(1,99) # 配置随机端口号
    id = str(uuid.uuid4()) #生成唯一id
	
    # 注册服务
    info = {'uuid':id,'ip':ip,'port':'{}'.format(port)}
    info = json.dumps(info) # 转为json
    logging.info('尝试注册 注册信息: {}'.format(info))
    r = requests.post(url='http://localhost:8888/regist', 
                      data=info,
                      headers={"Content-Type":"application/json"}
                      ) # 服务注册
    if r.status_code != 200:
        logging.info('注册服务失败 [{}]'.format(r.status_code))
        exit(-1)
    logging.info('注册服务成功 返回值:{}'.format(r.text))

if __name__ == "__main__":
    Init()
    logging.info('starting...')
    serve(app, host="0.0.0.0", port=port)
