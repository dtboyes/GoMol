### Required for all computers ###
To install required OpenGL, GLFW, and gonum libraries for graphics rendering/computation
```
go get -u github.com/go-gl/gl/v2.1/gl
go get -u github.com/go-gl/glfw/v3.3/glfw
go get -u gonum.org/v1/gonum/mat
```

### For Windows Machines ###
Getting the program set up for running on a Windows machine is a bit more involved.
1. Follow the installation instructions for [MSYS2](msys2.org).
   
   a. Required because it contains tools necessary for compiling and running native Windows programs.

   b. Contains the gcc compiler, which is necessary for running our program because the Go graphics libraries we are using are written in C or C++.
3. If MSYS2 is installed in the C Drive, run the following command to add the gcc executable to your PATH.
   
   a. ```$env:Path = "C:\msys64\ucrt64\bin;$env:Path"```

   b. Just replace directory with where msys64 directory is located on your machine.
   
5. Once this is complete, you can run the following command to open the web page, and follow all of the prompts to render your proteins of interest.
   
   a. ```python3 app.py```

   b. Ensure that you are in the correct directory when running this command.

### For Mac Machines ###
The process for getting set up on Mac machines is much simpler (after you have the required libraries).
1. Just run ```python3 app.py``` in the correct directory.
   a. The scene will only render in the bottom left quadrant of the window. This may be a system level issue on Mac, because the full window renders fine on Windows machines.
   


