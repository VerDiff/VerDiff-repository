import random

import matplotlib.pyplot as plt
import numpy as np

plt.rcParams["font.family"] = "serif"
plt.rcParams["font.serif"] = ["Times New Roman"] + plt.rcParams["font.serif"]

fig, ax = plt.subplots(figsize=(8, 2))


x = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36]

y1_file = open('./data/VerDiff_interesting_access.result', 'r')
y2_file = open('./data/Fluffy_interesting_access.result', 'r')

y1 = []
y2 = []

for line in y1_file.readlines():
    y1 = eval(line.strip())
for line in y2_file.readlines():
    y2 = eval(line.strip())


plt.plot(x, y1, label='VerDiff', marker='.', linestyle='-', color='blue',  markeredgecolor='blue')  
plt.plot(x, y2, label='Fluffy', marker='x', linestyle='-', color='green', markeredgecolor='green')

xtick_positions = list(range(3, 37, 3))
xtick_labels = [str(i) for i in range(1, 13)]
plt.xticks(xtick_positions, xtick_labels)

plt.yscale('log')

ax.set_ylim(10**4, 10**8)
ax.set_xlim(1, 36)

ax.set_ylabel('# of state accesses',fontsize=16)
ax.set_xlabel('Time (hours)',fontsize=16)
plt.xticks(fontsize=16)
plt.yticks(fontsize=16)

ax.legend(loc='upper left', fontsize=12, ncol=2)

plt.savefig('./diverse_access.pdf', dpi=600, bbox_inches='tight')
