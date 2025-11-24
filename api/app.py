from bottle import Bottle, run, template
import os
import weaviate

app = Bottle()

WEAVIATE_URL = os.environ["WEAVIATE_URL"]
client = weaviate.connect_to_local("weaviate", 8080)


@app.get('/')
def home():
    return {"message": "API is running"}


@app.route('/hello/<name>')
def index(name):
    return template('<b>Hello {{name}}</b>!', name=name)


@app.get('/weaviate/health')
def weaviate_health():
    return {
        "ready": client.is_ready(),
        "live": client.is_live(),
        "weaviate_url": WEAVIATE_URL
    }


run(app, host="0.0.0.0", port=8000)
