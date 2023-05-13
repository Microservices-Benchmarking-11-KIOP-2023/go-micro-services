from datetime import datetime
import json
import os
import matplotlib.pyplot as plt
import pandas as pd

test_output_dir = './test-scripts/test-outputs/'
metric_images_dir = './test-scripts/metric-images/'

filenames = os.listdir(test_output_dir)
filenames = [filename for filename in filenames if filename.endswith('.json')]
filenames = [os.path.join(test_output_dir, filename) for filename in filenames]

data_list = []
for filename in filenames:
    with open(filename) as f:
        data_list.append(json.load(f))


metrics = ['http_req_duration', 'http_req_blocked', 'http_req_connecting', 'http_req_waiting', 'http_req_receiving', 'http_req_sending', 'iteration_duration']
avg_metrics = {}

for metric in metrics:
    avg_metrics[metric] = []
    for data in data_list:
        try:
            avg_metrics[metric].append(data['metrics'][metric]['avg'])
        except KeyError:
            print(f"Key {metric} does not exist in one of the data files")
            avg_metrics[metric].append(None)

# Convert to DataFrame and set filenames as index
df = pd.DataFrame(avg_metrics, index=filenames)
# Directory for images from the current runtime
current_datetime = datetime.now()
datetime_str = current_datetime.strftime("%Y-%m-%d_%H-%M-%S")
os.mkdir(f'test-scripts\metric-images\\{datetime_str}')     

# Plot
for i, metric in enumerate(metrics, 1):
    df[metric].plot(kind='bar', ax=plt.gca())
    plt.title(f'Average {metric}')
    plt.xlabel('Filename')
    plt.ylabel('Duration [s]')
    plt.savefig(f'test-scripts\metric-images\{datetime_str}\{metric}.png',dpi=700,bbox_inches='tight',pad_inches=0.5)
