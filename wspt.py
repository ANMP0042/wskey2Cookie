import requests
import json
import os
import pyperclip
import base64
import logging
import time
import sys

requests.packages.urllib3.disable_warnings()
os.environ['no_proxy'] = '*'
logger = logging.getLogger(__name__)

def cloud_info(url_t):
    url = str(base64.b64decode(url_t).decode()) + 'check_api'
    for i in range(3):
        try:
            headers = {"authorization": "Bearer Shizuku"}
            res = requests.get(url=url, verify=False, headers=headers, timeout=20).text
        except requests.exceptions.ConnectTimeout:
            logger.info("\n获取云端参数超时, 正在重试!" + str(i))
            time.sleep(1)
            continue
        except requests.exceptions.ReadTimeout:
            logger.info("\n获取云端参数超时, 正在重试!" + str(i))
            time.sleep(1)
            continue
        except Exception as err:
            logger.info("\n未知错误云端, 退出脚本!")
            logger.debug(str(err))
            sys.exit(1)
        else:
            try:
                c_info = json.loads(res)
            except Exception as err:
                logger.info("云端参数解析失败")
                logger.debug(str(err))
                sys.exit(1)
            else:
                return c_info


def check_cloud():
    url_list = ['aHR0cDovLzQzLjEzNS45MC4yMy8=', 'aHR0cHM6Ly9zaGl6dWt1Lm1sLw==', 'aHR0cHM6Ly9jZi5zaGl6dWt1Lm1sLw==']
    for i in url_list:
        url = str(base64.b64decode(i).decode())
        try:
            requests.get(url=url, verify=False, timeout=10)
        except Exception as err:
            logger.debug(str(err))
            continue
        else:
            info = ['Default', 'HTTPS', 'CloudFlare']
            logger.info(str(info[url_list.index(i)]) + " Server Check OK\n--------------------\n")
            return i
    logger.info("\n云端地址全部失效, 请检查网络!")
    sys.exit(1)

# 返回值 bool jd_ck
def getToken(wskey):
    url_t = check_cloud()
    ua = cloud_info(url_t)['User-Agent']
    try:
        url = str(base64.b64decode(url_t).decode()) + 'genToken'
        header = {"User-Agent": ua}
        params = requests.get(url=url, headers=header, verify=False, timeout=20).json()
    except Exception as err:
        logger.info("Params参数获取失败")
        logger.debug(str(err))
        return False, wskey
    headers = {
        'cookie': wskey,
        'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'charset': 'UTF-8',
        'accept-encoding': 'br,gzip,deflate',
        'user-agent': ua
    }
    url = 'https://api.m.jd.com/client.action'
    data = 'body=%7B%22to%22%3A%22https%253a%252f%252fplogin.m.jd.com%252fjd-mlogin%252fstatic%252fhtml%252fappjmp_blank.html%22%7D&'
    try:
        res = requests.post(url=url, params=params, headers=headers, data=data, verify=False, timeout=10)
        res_json = json.loads(res.text)
        tokenKey = res_json['tokenKey']
    except Exception as err:
        logger.info("JD_WSKEY接口抛出错误 尝试重试 更换IP")
        logger.info(str(err))
        return False, wskey
    else:
        return appjmp(wskey, tokenKey, ua)


# 返回值 bool jd_ck
def appjmp(wskey, tokenKey, ua):
    wskey = "pt_" + str(wskey.split(";")[0])
    if tokenKey == 'xxx':
        logger.info(str(wskey) + ";WsKey状态失效\n--------------------\n")
        return False, wskey
    headers = {
        'User-Agent': ua,
        'accept': 'accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
        'x-requested-with': 'com.jingdong.app.mall'
    }
    params = {
        'tokenKey': tokenKey,
        'to': 'https://plogin.m.jd.com/jd-mlogin/static/html/appjmp_blank.html',
    }
    url = 'https://un.m.jd.com/cgi-bin/app/appjmp'
    try:
        res = requests.get(url=url, headers=headers, params=params, verify=False, allow_redirects=False, timeout=20)
    except Exception as err:
        logger.info("JD_appjmp 接口错误 请重试或者更换IP\n")
        logger.info(str(err))
        return False, wskey
    else:
        try:
            res_set = res.cookies.get_dict()
            pt_key = 'pt_key=' + res_set['pt_key']
            pt_pin = 'pt_pin=' + res_set['pt_pin']
            jd_ck = str(pt_key) + '; ' + str(pt_pin) + '; __time=' + str(time.time())
        except Exception as err:
            logger.info("JD_appjmp提取Cookie错误 请重试或者更换IP\n")
            logger.info(str(err))
            return False, wskey
        else:
            if 'fake' in pt_key:
                logger.info(str(wskey) + ";WsKey状态失效\n")
                return False, wskey
            else:
                logger.info(str(wskey) + ";WsKey状态正常\n")
                return True, jd_ck


if __name__ == '__main__':
    os.system('chcp 65001')
    print("Ver: 0.1")
    print("limoe")
    #print("请输入格式为pin=xxx;wskey=xxx;")
    ws = input()
    url_t = check_cloud()
    args = cloud_info()
    ua = args['User-Agent']
    return_ws = getToken(ws)
    print(return_ws)
    os.system('pause')
