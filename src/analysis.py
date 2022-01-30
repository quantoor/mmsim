import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from datetime import datetime


font = {
        # 'family': 'serif',
        # 'color':  'darkred',
        'weight': 'normal',
        'size': 12,
        }


def plotMetrics(df, figName, savePath=None):
    fig, (ax1, ax2, ax3, ax4) = plt.subplots(4, 1, figsize=(16, 8), sharex='col')
    fig.suptitle(figName, fontsize=14)

    x = df.index
    p0 = df['Price'][0]
    e0 = df['Balance'][0]

    ax1.set_title('Price')
    ax1.plot(x, np.array(df['Price']), linewidth=1.3)
    lbl = ax1.set_ylabel('$', labelpad=10)
    lbl.set_rotation(0)
    ax1b = ax1.twinx()
    ax1b.plot(x, (df['Price'] - p0) / p0 * 100, linewidth=0)
    lbl = ax1b.set_ylabel('%', labelpad=10)
    lbl.set_rotation(0)

    ax2.set_title('Balance')
    ax2.plot(x, df['Balance'], linewidth=0.5)
    ax2.fill_between(x, df['Balance'][0], df['Balance'], alpha=0.5)
    lbl = ax2.set_ylabel('$', labelpad=10)
    lbl.set_rotation(0)
    ax2.ticklabel_format(axis='y', useOffset=False)
    ax2b = ax2.twinx()
    ax2b.plot(x, (df['Balance'] - e0) / e0 * 100, linewidth=0)
    lbl = ax2b.set_ylabel('%', labelpad=10)
    lbl.set_rotation(0)

    # pnl = 'PNL-' + side
    # ax3.set_title('Unrealized PNL')
    # ax3.plot(x, df[pnl], linewidth=0)
    # ax3.fill_between(x, df[pnl][0], df[pnl], where=df[pnl]>0, alpha=0.5, color='green')
    # ax3.fill_between(x, df[pnl][0], df[pnl], where=df[pnl]<0, alpha=0.5, color='red')
    # lbl = ax3.set_ylabel('$', labelpad=10)
    # lbl.set_rotation(0)
    # ax3.ticklabel_format(axis='y', useOffset=False)
    # ax3b = ax3.twinx()
    # ax3b.plot(x, df[pnl] / df['Equity'] * 100, linewidth=0)
    # lbl = ax3b.set_ylabel('%', labelpad=10)
    # lbl.set_rotation(0)

    # ax4.set_title('Equity + Unrealized PNL')
    # df['equity+pnl'] = df['Equity'] + df[pnl]
    # ax4.plot(x, df['equity+pnl'], linewidth=0.5)
    # ax4.fill_between(x, df['Equity'][0], df['equity+pnl'], alpha=0.5)
    # lbl = ax4.set_ylabel('$', labelpad=10)
    # lbl.set_rotation(0)
    # ax4.ticklabel_format(axis='y', useOffset=False)
    # ax4b = ax4.twinx()
    # ax4b.plot(x, (df['equity+pnl'] - e0) / e0 * 100, linewidth=0)
    # lbl = ax4b.set_ylabel('%', labelpad=10)
    # lbl.set_rotation(0)
 
    fig.tight_layout()
    if savePath is not None:
        fig.savefig(savePath)


def plotDistributions(df, side, savePath=None):
    profit = 'NetProfit-' + side
    gridReached = 'GridReached-'+ side

    df[profit] = df['GrossProfit-L']
    fig, (ax1, ax2, ax3) = plt.subplots(1, 3, figsize=(16, 5))
    # fig.suptitle('Distribution of take profits', fontsize=14)
        
    ### get grid reached for each profit
    df_profit = df.loc[(df[profit] != 0)]
    profitGridCount = df_profit[gridReached].value_counts()

    data = np.array(profitGridCount)
    ax1.bar(profitGridCount.index, data)
    ax1.set_title('Distribution of TP grids', fontdict=font)
    ax1.set_xlabel('Grid Number', fontdict=font)

    ax1.set_ylabel('Occurrences', fontdict=font)
    ax1b = ax1.twinx()
    ax1b.bar(profitGridCount.index, data/data.sum() * 100)
    ax1b.set_ylabel('%', labelpad=10, fontdict=font)

    ### compute profit contribution of each grid
    df_profit_distr = df_profit.groupby([gridReached]).sum()
    df_profit_distr = df_profit_distr[profit]

    data = np.array(df_profit_distr)
    ax2.bar(df_profit_distr.index, data)
    ax2.set_title('Distribution of profit per grid', fontdict=font)
    ax2.set_xlabel('Grid Number', fontdict=font)

    ax2.set_ylabel('$', labelpad=10, fontdict=font)
    ax2b = ax2.twinx()
    ax2b.bar(df_profit_distr.index, data / abs(data).sum() * 100)
    ax2b.set_ylabel('%', labelpad=10, fontdict=font)    

    fig.tight_layout()
    fig.subplots_adjust(wspace=0.4)
    if savePath is not None:
        fig.savefig(savePath)


    ### mean profit per grid
    df_profit_distr = df_profit.groupby([gridReached]).mean()
    df_profit_distr = df_profit_distr[profit]
    ax3.bar(df_profit_distr.index, np.array(df_profit_distr))
    ax3.set_title('Mean profit per grid', fontdict=font)
    ax3.set_xlabel('Grid Number', fontdict=font)
    ax3.set_ylabel('$', labelpad=10, fontdict=font)

    return
    ### consecutive tp at grid 0
    df = df.loc[(df['GrossProfit-L'] != 0)]
    df['subgroup'] = (df['GridReached-L'] != df['GridReached-L'].shift(1)).cumsum()
    df = df[df['GridReached-L'] == 0] # drop all rows where grid reached is not 0
    df = df[['Timestamp', 'GridReached-L', 'subgroup']]    
    df = df.groupby('subgroup',as_index=False).apply(f)
    df = df[df['N'] > 1] # remove rows with N==1
    # print(df.head(10))
    # print(df['N'].value_counts())
    # print(df['N'].sum())
    df = df['N'].value_counts()
    print(df)
    # data = np.array(df)
    # ax4.bar(df.index, data)
    # ax4.set_title('Distribution of number of consecutive TP at grid 0', fontdict=font)
    # ax4.set_xlabel('Consecutive TP', fontdict=font)
    # ax4.set_ylabel('Occurences', labelpad=10, fontdict=font)


def f(x):
    N = x.shape[0]
    dt = x['Timestamp'].iloc[-1] - x['Timestamp'].iloc[0]
    return pd.Series([N, dt/N], index=['N', 'mean_time_s'])


if __name__ == '__main__':
    folder = "results/"
    file = "doge_results.csv"
    df = pd.read_csv(folder + file)
    df.index = [datetime.fromtimestamp(x) for x in df['Timestamp']]

    plotMetrics(df, file)
    # plotDistributions(df)
    plt.show()