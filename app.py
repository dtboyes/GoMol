from flask import Flask, render_template, request
import os
import subprocess

app = Flask(__name__)

UPLOAD_FOLDER = r"C:\Users\BLUE OCEAN\Desktop\FlaskWeb\GoMol-main\GoMol-main"
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

BASE_DIR = os.path.abspath(os.path.dirname(__file__))

@app.route('/')
@app.route('/home')
def home():
    return render_template("index.html")

@app.route('/result', methods=['POST', 'GET'])
def result():
    
    if request.method == 'POST':
        pdb_id_1 = request.form.get('pdb_id_1')
        pdb_id_2 = request.form.get('pdb_id_2')
        render_chain_a = request.form.get('render_chain_a')

        # Save data to a text file
        save_to_text(pdb_id_1, pdb_id_2, render_chain_a)

        return render_template('result.html')

def save_to_text(pdb_id_1, pdb_id_2, render_chain_a):
    # Define the path to the text file using os.path.join
    text_file_path = os.path.join(app.config['UPLOAD_FOLDER'], 'user_data.txt')

    # Save the data to the text file
    with open(text_file_path, 'a') as text_file:
        text_file.write(f"{pdb_id_1}\n")
        text_file.write(f"{pdb_id_2}\n")
        text_file.write(f"{render_chain_a}\n")
        text_file.write("\n")

@app.route('/result2')
def result2():
    # Fetch the data from the text file
    text_file_path = os.path.join(app.config['UPLOAD_FOLDER'], 'user_data.txt')
    with open(text_file_path, 'r') as text_file:
        data = text_file.read()

    # Execute the GoMol program (assuming data is formatted as expected)
    execute_gomol(*data.split('\n')[0:3])  # assuming data has PDB ID1, PDB ID2, Render Chain A in the first 3 lines

    return render_template('result2.html')

def execute_gomol(pdb_id_1, pdb_id_2, render_chain_a):
    # Execute the GoMol program and pass the user data
    subprocess.run([r"C:\Users\BLUE OCEAN\Desktop\FlaskWeb\GoMol-main\GoMol-main\GoMol\GoMol.exe", pdb_id_1, pdb_id_2, render_chain_a])

if __name__ == "__main__":
    app.run(debug=True)