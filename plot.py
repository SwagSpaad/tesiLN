import pandas as pd 
import matplotlib.pyplot as plt
import os

plt.style.use('seaborn-v0_8-paper')
plt.rcParams.update({
    'font.size': 16,
    'axes.labelsize': 16,
    'axes.titlesize': 16,
    'legend.fontsize': 16,
    'xtick.labelsize': 14,
    'ytick.labelsize': 14,
    'lines.linewidth': 2.5,
    'lines.markersize': 6,
    'figure.figsize': (14, 8),
})

def load_data(filename):
    if not os.path.exists(filename):
        print(f"File {filename} non trovato.")
        return None
    
    df = pd.read_csv(filename)
    return {
        'x': df.iloc[:, 0],
        'liqFail': df.iloc[:, 1],
        'pathFail': df.iloc[:, 2],
        'avgHops': df.iloc[:, 3]
    }

print("Carico dati CSV...")
data_LN_casuali = load_data('LN_GuastiCasuali_noGiant.csv')
data_LN_mirati = load_data('LN_AttacchiMirati_noGiant.csv')
data_ER_casuali = load_data('ER_GuastiCasuali_noGiant.csv')
data_ER_mirati = load_data('ER_AttacchiMirati_noGiant.csv')

colore_ln = '#d62728'
colore_er = '#1f77b4'

if all([data_LN_casuali, data_LN_mirati, data_ER_casuali, data_ER_mirati]):
    #Grafico 1 guasti casuali - fallimento path
    plt.figure()
    plt.plot(data_LN_casuali['x'], data_LN_casuali['pathFail'], color=colore_ln, linestyle='-', label='Lightning Network')
    plt.plot(data_ER_casuali['x'], data_ER_casuali['pathFail'], color=colore_er, linestyle='-', label='Erdos-Renyi')
    plt.xlabel('Nodi rimossi')
    plt.ylabel('Fallimenti per assenza percorso (%)')
    plt.xlim(0, 150)
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend()
    plt.tight_layout()
    plt.savefig('1_GuastiCasuali_path_noGiant.pdf')
    print('Salvato: 1_GuastiCasuali_path.pdf')

    # Grafico 2 guasti casuali - fallimento liquidità
    plt.figure()
    plt.plot(data_LN_casuali['x'], data_LN_casuali['liqFail'], color=colore_ln, linestyle='-', label='Lightning Network')
    plt.plot(data_ER_casuali['x'], data_ER_casuali['liqFail'], color=colore_er, linestyle='-', label='Erdos-Renyi')
    plt.xlabel('Nodi rimossi')
    plt.ylabel('Fallimenti per liquidità insufficiente (%)')
    plt.xlim(0, 150)
    plt.ylim(0, 100)
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend()
    plt.tight_layout()
    plt.savefig('2_GuastiCasuali_liq_noGiant.pdf')
    print('Salvato: 2_GuastiCasuali_liq.pdf')

    #Grafico 3 attacchi mirati - fallimento path
    plt.figure()
    plt.plot(data_LN_mirati['x'], data_LN_mirati['pathFail'], color=colore_ln, linestyle='-', linewidth=2.5, label='Lightning Network')
    plt.plot(data_ER_mirati['x'], data_ER_mirati['pathFail'], color=colore_er, linestyle='-', label='Erdos-Renyi')
    plt.fill_between(data_LN_mirati['x'], data_LN_mirati['pathFail'], alpha=0.1, color=colore_ln)
    plt.xlabel('Nodi rimossi')
    plt.ylabel('Fallimenti per assenza percorso (%)')
    plt.xlim(0, 150)
    plt.ylim(0, 100)
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend()
    plt.tight_layout()
    plt.savefig('3_AttracchiMirati_path_noGiant.pdf')
    print('Salvato: 3_AttacchiMirati_path.pdf')

    #Grafico 4 attacchi mirati - fallimento liquidità
    plt.figure()
    plt.plot(data_LN_mirati['x'], data_LN_mirati['liqFail'], color=colore_ln, linestyle='-', linewidth=2.5, label='Lightning Network')
    plt.plot(data_ER_mirati['x'], data_ER_mirati['liqFail'], color=colore_er, linestyle='-', label='Erdos-Renyi')
    plt.xlabel('Nodi rimossi')
    plt.ylabel('Fallimenti per liquidità insufficiente (%)')
    plt.xlim(0, 150)
    plt.ylim(0, 100)
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend()
    plt.tight_layout()
    plt.savefig('4_AttracchiMirati_liq_noGiant.pdf')
    print('Salvato: 4_AttacchiMirati_liq.pdf')

    #Grafico 5 analisi aumento hop medi
    plt.figure()
    plt.plot(data_LN_mirati['x'], data_LN_mirati['avgHops'], color=colore_ln, linestyle='-', label='Lightning Network (AM)')
    plt.plot(data_ER_mirati['x'], data_ER_mirati['avgHops'], color=colore_er, linestyle='-', label='Erdos-Renyi (AM)')
    plt.plot(data_LN_casuali['x'], data_LN_casuali['avgHops'], color=colore_ln, linestyle='--', alpha=0.5, label='Lightning Network (GC)')
    plt.plot(data_ER_casuali['x'], data_ER_casuali['avgHops'], color=colore_er, linestyle='--', alpha=0.5, markevery=5, label='Erdos-Renyi (GC)')

    plt.xlabel('Nodi Rimossi')
    plt.ylabel('Numero medio Hop')
    plt.xlim(0, 150)
    plt.ylim(3.5, 10)
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend(loc='upper left', fontsize=10)
    plt.tight_layout()
    plt.savefig('5_AumentoHop_noGiant.pdf')
    print("Salvato: 5_AumentoHop.pdf")

print("Grafici generati con successo")