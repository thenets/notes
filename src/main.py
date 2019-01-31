import hashlib, redis, datetime, humanize, sys

from flask import Flask, render_template, url_for, request, jsonify
app = Flask(__name__)

# Connect to the Redis
r = redis.StrictRedis(
    host='localhost',
    port=6379)


def getNoteByPath (path):
      """Returns the note based on path

      @params path <string>
      @return dict[content, updateAt, updateAtHuminze, error]
      """
      out = {}
      out['content'] = ''
      out['updateAt'] = ''
      out['updateAtHumanize'] = ''
      out['error'] = ''

      # Check if note exist
      if r.exists(path) == 0:
            out['error'] = "Note doesn't exist!"
      
      else:
            # Get note from the database
            try:
                  updateAt = r.get(path+':time').decode('utf-8')
                  updateAtHumanize = humanize.naturaltime( datetime.datetime.now() - datetime.datetime.strptime(updateAt, '%Y-%m-%d %H:%M:%S') )

                  out['content'] = r.get(path).decode('utf-8')
                  out['updateAt'] = updateAt
                  out['updateAtHumanize'] = updateAtHumanize
            except:
                  out['error'] = "Fail to read the note!"

      return out




###
# Notes API
###

# Read note from the database
@app.route('/api/<path:path>', methods=['GET'])
def api_read(path):
      out = getNoteByPath(path)

      return jsonify(out)

# Write note to the database
@app.route('/api/<path:path>', methods=['PUT', 'POST'])
def api_write(path):
      out = {}
      out['content'] = ''
      out['updateAt'] = ''
      out['error'] = ''

      # Check if note exist
      try:
            out['content'] = request.form['note']
      except:
            out['error'] = "The 'note' key doesn't exist!"

      # Store the note in the database
      try:
            updateAt = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

            # Delete if 'note' is empty
            if len(request.form['note']) == 0:
                  r.delete(path)
                  r.delete(path+':time')
            else:
                  r.set(path, request.form['note'])
                  r.set(path+':time', updateAt)
      except Exception as e:
            out['error'] = "Fail to save the note!"
            print(str(e))
      
      out = getNoteByPath(path)

      return jsonify(out)


##
# View render
##

@app.route('/favicon.ico', methods=['GET'])
def favicon():
      return ''

# All notes requests
@app.route('/<path:path>', methods=['GET'])
def catch_all(path):
      out = getNoteByPath(path)

      # Render note view
      return render_template('note.html',
                              path=path,
                              content=out['content'],
                              updateAt=out['updateAt'],
                              updateAtHumanize=out['updateAtHumanize'],
                              error=out['error'])
    
    
# Default homepage
@app.route('/')
def hello_world():
      return render_template('home.html')

if __name__ == '__main__':
      app.run(host='0.0.0.0', debug=True, port=5000)
