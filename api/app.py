from bottle import Bottle, run, template

app = Bottle()

@app.get('/')
def home():
    return {"message": "API is running"}

@app.route('/hello/<name>')
def index(name):
    return template('<b>Hello {{name}}</b>!', name=name)

run(app, host="0.0.0.0", port=8000)