import sqlite3
import datetime
import random

# Connect to SQLite database (or create it if it doesn't exist)
conn = sqlite3.connect('health_checks.db')
cursor = conn.cursor()

# Create the health_check_results table if it doesn't exist
cursor.execute('''
CREATE TABLE IF NOT EXISTS health_check_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_endpoint TEXT NOT NULL,
    http_method TEXT NOT NULL,
    callout_timestamp TIMESTAMP NOT NULL,
    response_code INTEGER NOT NULL,
    response TEXT,
    time_taken FLOAT NOT NULL
)
''')
conn.commit()

# Define the 4 target endpoints
endpoints = ['https://api.anthropic.com/v1/messages', 'https://api.openai.com/v1/chat/completions', 'https://google.com', 'https://citi.com']

# Set the start and end times
start_time = datetime.datetime.now() - datetime.timedelta(days=60)
end_time = datetime.datetime.now()

# Initialize current_time to start_time
current_time = start_time

# Batch settings to optimize database insertion
batch_size = 1000
data_batch = []

# Loop over each 10-minute interval
while current_time <= end_time:
    for endpoint in endpoints:
        target_endpoint = endpoint
        http_method = 'GET'  # You can randomize this if needed
        callout_timestamp = current_time

        # Generate response_code with 80% chance of being 200
        response_code = random.choices(
            population=[200, 404, 500, 503],
            weights=[80, 10, 5, 5],
            k=1
        )[0]

        # Generate corresponding response text
        if response_code == 200:
            response = 'OK'
        elif response_code == 404:
            response = 'Not Found'
        elif response_code == 500:
            response = 'Internal Server Error'
        elif response_code == 503:
            response = 'Service Unavailable'
        else:
            response = 'Error'

        # Simulate time taken between 0.1 and 1.0 seconds
        time_taken = round(random.uniform(0.1, 1.0), 3)

        # Append the generated data to the batch
        data_batch.append((
            target_endpoint,
            http_method,
            callout_timestamp,
            response_code,
            response,
            time_taken
        ))

        # Insert data in batches to optimize performance
        if len(data_batch) >= batch_size:
            cursor.executemany('''
            INSERT INTO health_check_results (
                target_endpoint,
                http_method,
                callout_timestamp,
                response_code,
                response,
                time_taken
            ) VALUES (?, ?, ?, ?, ?, ?)
            ''', data_batch)
            conn.commit()
            data_batch = []  # Reset the batch

    # Increment the current time by 10 minutes
    current_time += datetime.timedelta(minutes=10)

# Insert any remaining data in the batch
if data_batch:
    cursor.executemany('''
    INSERT INTO health_check_results (
        target_endpoint,
        http_method,
        callout_timestamp,
        response_code,
        response,
        time_taken
    ) VALUES (?, ?, ?, ?, ?, ?)
    ''', data_batch)
    conn.commit()

# Close the database connection
conn.close()
