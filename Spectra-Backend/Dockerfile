FROM python:3.10

WORKDIR /code

# Install system dependencies required by OpenCV and Ultralytics
RUN apt-get update && apt-get install -y libgl1 libglib2.0-0

COPY ./requirements.txt /code/requirements.txt

# Install the Python dependencies
RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

# Copy your main.py file into the container
COPY . .

# Run the FastAPI server on port 7860
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "7860"]