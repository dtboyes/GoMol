from flask import Flask, render_template, request
import os

app = Flask(__name__)

UPLOAD_FOLDER = 'C:/Users/BLUE OCEAN/Desktop/FlaskApp/GoMol-main/GoMol/pdbfiles'
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

@app.route('/')
@app.route('/home')
def home():
    return render_template("index.html")

@app.route('/result', methods=['POST', 'GET'])
def result():
    if request.method == 'POST':
        # Get form data
        name = request.form.get('name')

        # Check if the post request has the file part
        if 'file' in request.files:
            file = request.files['file']
            
            # Save the file to the specified directory
            if file.filename != '':
                file_path = os.path.join(app.config['UPLOAD_FOLDER'], file.filename)
                file.save(file_path)
            else:
                file_path = None
        else:
            file_path = None

        # Pass data to the template
        return render_template('result.html', name=name, file_path=file_path)

@app.route('/result2')
def result2():
    # Get data from the query parameters
    name = request.args.get('name')
    file_path = request.args.get('file_path')

    # Pass data to the template
    return render_template('result2.html', name=name, file_path=file_path)

if __name__ == "__main__":
    app.run(debug=True)