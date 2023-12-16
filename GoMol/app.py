from flask import Flask, render_template, request
import os
import subprocess

app = Flask(__name__)


BASE_DIR = os.path.abspath(os.path.dirname(__file__))
app.config['UPLOAD_FOLDER'] = BASE_DIR

  
@app.route('/')
@app.route('/home')
def home():
    return render_template("index.html")

@app.route('/result', methods=['POST', 'GET'])
def result():
    
    if request.method == 'POST':
        if os.path.exists("user_data.txt"):
            os.remove("user_data.txt")
        pdb_id_1 = request.form.get('pdb_id_1')
        pdb_id_2 = request.form.get('pdb_id_2')

        # Save data to a text file
        save_to_text(pdb_id_1, pdb_id_2)

        return render_template('result.html')

def save_to_text(pdb_id_1, pdb_id_2):
    # Define the path to the text file using os.path.join
    text_file_path = os.path.join(app.config['UPLOAD_FOLDER'], 'user_data.txt')

    # Save the data to the text file
    with open(text_file_path, 'a') as text_file:
        text_file.write(f"{pdb_id_1}\n")
        text_file.write(f"{pdb_id_2}\n")
        text_file.write("\n")

@app.route('/result2')
def result2():
    # Fetch the data from the text file
    text_file_path = os.path.join(app.config['UPLOAD_FOLDER'], 'user_data.txt')
    with open(text_file_path, 'r') as text_file:
        data = text_file.read()

    # Execute the GoMol program (assuming data is formatted as expected)
    execute_gomol(*data.split('\n')[0:2])  # assuming data has PDB ID1, PDB ID2, Render Chain A in the first 3 lines
    line1, line2, line3, line4 = display_sequence()

    # Split the string into individual score strings
    score_strings = line4.split()

    # Convert each score string to a float and store it in an array
    qres_scores = [float(score) for score in score_strings]

    colors = []
    score_index = 0
    for char in line2:
        if char == '|':
            color = score_to_color(qres_scores[score_index])
            score_index += 1
        else:
            color = 'rgb(255, 0, 0)'  # Red for unaligned residues
        colors.append(color)
    zipped_data = zip(line2, colors)
    return render_template('result2.html', seq1= line1, match=line2, seq2=line3, zipped=zipped_data)

def execute_gomol(pdb_id_1, pdb_id_2):
    # Execute the GoMol program and pass the user data
    if os.name == "nt":
        path = os.path.join(BASE_DIR, "GoMol.exe")
    elif os.name == "posix":
        path = os.path.join(BASE_DIR, "GoMol")
    subprocess.run([path, pdb_id_1, pdb_id_2])


def score_to_color(score):
    # Convert score to color: higher score -> more blue, lower score -> more red
    red = int(255 * (1 - score))
    blue = int(255 * score)
    return f'rgb({red}, 0, {blue})'



def display_sequence():
    # sequence.txt 파일 읽기
    file_path = 'result.txt'

    with open(file_path, 'r') as file:
        lines = file.readlines()

    line1 = lines[0]
    line2 = lines[1]
    line3 = lines[2]
    line4 = lines[3]



    return line1, line2, line3, line4


def read_lines(filename):
    # 파일 경로
    file_path = os.path.join(os.path.dirname(__file__), filename)
    
    # 파일 읽기
    with open(file_path, 'r') as file:
        lines = file.readlines()
    
    return lines

if __name__ == "__main__":
    app.run(debug=True)