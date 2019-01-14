import hashlib, redis, datetime

from flask import Flask, render_template, url_for, request, jsonify
app = Flask(__name__)

# Connect to the Redis
r = redis.StrictRedis(
    host='localhost',
    port=6379)

###
# Notes API
###

# Read note from the database
@app.route('/api/<path:path>', methods=['GET'])
def api_read(path):
      out = {}

      # Check if note exist
      if r.exists(path) == 0:
            out['error'] = "Note doesn't exist!"
      
      else:
            # Get note from the database
            try:
                  out['content'] = r.get(path).decode('utf-8')
                  out['updateAt'] = r.get(path+':time').decode('utf-8')
            except:
                  out['error'] = "Fail to read the note!"

      return jsonify(out)

# Write note to the database
@app.route('/api/<path:path>', methods=['PUT', 'POST'])
def api_write(path):
      out = {}

      # Check if note exist
      try:
            out['content'] = request.form['note']
      except:
            out['error'] = "The 'note' key doesn't exist!"

      # Store the note in the database
      try:
            r.set(path, request.form['note'])
            r.set(path+':time', datetime.datetime.now().isoformat())
            out['content'] = r.get(path).decode('utf-8')
            out['updateAt'] = r.get(path+':time').decode('utf-8')
      except:
            out['error'] = "Fail to save the note!"

      return jsonify(out)


##
# View render
##

# All notes requests
@app.route('/<path:path>', methods=['GET'])
def catch_all(path):

      # Note update
      if request.method == 'PUT':
            return "cafe"
      
      # Get note content by hash

      # Note render view
      return 'You want path: %s' % path


    
# Default homepage
@app.route('/')
def hello_world():
      return render_template('home.html')

if __name__ == '__main__':
      app.run(host='0.0.0.0', debug=True, port=5000)