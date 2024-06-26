from dotenv import load_dotenv
import os

load_dotenv()  

HTTP_NODE_URL = os.getenv('HTTP_NODE_URL')
WEBSOCKET_NODE_URL = os.getenv('WEBSOCKET_NODE_URL')