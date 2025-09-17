import matplotlib.pyplot as plt
import matplotlib.ticker as ticker
from icecream import ic

plt.rcParams["font.family"] = "serif"
plt.rcParams["font.serif"] = ["Times New Roman"] + plt.rcParams["font.serif"]


def num_of_suicide_contracts():
    fig = plt.figure(figsize=(8, 2))

    x = []
    y = []
    count = 0
    with open('./data/Num_selfdestruct_contracts') as f:
        for line in f:
            count += 1
            if count == 1:
                 continue
            info = line.strip().split(',')
            x.append(int(info[0]))
            y.append(int(info[1]))

    plt.plot(x, y, "rs-", label="30", markersize=12, markerfacecolor="none")
    plt.yscale('log')
    #plt.xscale('log')
    
    plt.xlim(50000, 20000000)
    plt.ylim(1, 10**9)

    plt.xticks([50000,2500000,5000000,7500000,10000000,12500000,15000000,17500000, 20000000], labels=['0M','2.5M','5M','7.5M','10M','12.5M','15M', '17.5M', '20M'])
    plt.xticks(rotation=30)
    plt.annotate(
        '27,372,264',
        xy=(19750001,27372264),
        xytext=(16000000, 100000),
        fontsize=14,
        color="k",
        ha='center',
        va='center',
        arrowprops=dict(arrowstyle="->", connectionstyle="arc3", color="k"),
    )

    plt.xlabel("Block height", fontsize=16)
    plt.ylabel("# of contracts", fontsize=16)

    plt.xticks(fontsize=15)
    plt.yticks(fontsize=15)

    plt.savefig("./num_of_selfdestruct_contracts.pdf", dpi=1000, bbox_inches="tight")
    
num_of_suicide_contracts()
