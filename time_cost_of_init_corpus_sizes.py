import matplotlib.pyplot as plt
import numpy as np

plt.rcParams["font.family"] = "serif"
plt.rcParams["font.serif"] = ["Times New Roman"] + plt.rcParams["font.serif"]

fig, ax = plt.subplots(figsize=(8, 2.5))

x1 = [0, 5, 5, 9, 9, 16, 16, 34, 34, 55, 55, 93, 93, 107, 107, 154, 154, 200] # size = 100
x2 = [0, 10, 10, 15, 15, 23, 23, 44, 44, 69, 69, 103, 103, 135, 135, 162, 162, 200] # size = 1000
x3 = [0, 18, 18, 29, 29, 38, 38, 57, 57, 107, 107, 133, 133, 152, 152, 191, 191, 200] # size = 10000
y = [0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8]

plt.plot(x1, y, label=r'$\mathit{Corpus}_{100}$', linestyle='-', color='black')    
plt.plot(x2, y, label=r'$\mathit{Corpus}_{1000}$', linestyle='--', color='black')  
plt.plot(x3, y, label=r'$\mathit{Corpus}_{10000}$', linestyle='-.', color='black')  

ax.set_ylim(0, 9)
ax.set_xlim(0, 200)

ax.set_ylabel('# of identified issues',fontsize=16)
ax.set_xlabel('Time (minutes)',fontsize=16)
plt.xticks(fontsize=16)
plt.yticks(fontsize=16)

ax.legend(loc='upper left', fontsize=12, ncol=1)

plt.savefig('./time_cost.pdf', dpi=600, bbox_inches='tight')
