import pandas as pd
import matplotlib.pyplot as plt
import os
import numpy as np

# Configurazione stile accademico
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

# Nomi dei file
file_ln = 'GradoLN.csv'
file_er = 'GradoER.csv'

# Colori per coerenza con i grafici precedenti
color_ln = '#d62728'  # Rosso
color_er = '#1f77b4'  # Blu

print("Caricamento dati CSV...")
df_ln = pd.read_csv(file_ln) if os.path.exists(file_ln) else None
df_er = pd.read_csv(file_er) if os.path.exists(file_er) else None

if df_ln is not None and df_er is not None:
    
    # =========================================================================
    # GRAFICO 1: Crescita del grado (Dati Ordinati)
    # =========================================================================
    plt.figure()
    
    # Asse X: da 0 al numero totale di nodi
    # Asse Y: il grado del nodo
    plt.plot(range(len(df_ln)), df_ln['Grado'], color=color_ln, label='Lightning Network')
    plt.plot(range(len(df_er)), df_er['Grado'], color=color_er, label='Erdős-Rényi')
    
    plt.xlabel('Indice del Nodo')
    plt.ylabel('Grado del Nodo (Numero di Canali)')
    plt.yscale('log') # Usiamo la scala logaritmica sulla Y per non schiacciare troppo i nodi ER
    
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend()
    plt.tight_layout()
    plt.savefig('6_Crescita_Gradi_Ordinati.pdf')
    print("Salvato: 6_Crescita_Gradi_Ordinati.pdf")


    # =========================================================================
    # GRAFICO 2: Distribuzione di Probabilità P(k) in Scala Log-Log
    # =========================================================================
    plt.figure()

    # Funzione per calcolare la probabilità P(k) = N_k / N_tot
    def get_probability_distribution(df):
        # Conta quanti nodi hanno un determinato grado
        degree_counts = df['Grado'].value_counts().sort_index()
        # Normalizza per ottenere la probabilità
        probabilities = degree_counts / len(df)
        return probabilities.index, probabilities.values

    degrees_ln, probs_ln = get_probability_distribution(df_ln)
    degrees_er, probs_er = get_probability_distribution(df_er)

    # Disegniamo i punti in scala Log-Log
    plt.loglog(degrees_ln, probs_ln, marker='o', linestyle='none', color=color_ln, alpha=0.7, label='Lightning Network')
    #plt.loglog(degrees_er, probs_er, marker='s', linestyle='none', color=color_er, alpha=0.7, label='Erdős-Rényi (Poisson)')

    plt.xlabel('Grado $k$')
    plt.ylabel('Probabilità $P(k)$')
    
    plt.grid(True, which="both", linestyle=':', alpha=0.5)
    plt.legend()
    plt.tight_layout()
    plt.savefig('7_Distribuzione_Gradi_LogLog.pdf')
    print("Salvato: 7_Distribuzione_Gradi_LogLog.pdf")

    # =========================================================================
    # GRAFICO 3: Distribuzione P(k) in Scala Lineare (La Campana di ER)
    # =========================================================================
    plt.figure()

    # Disegniamo i punti in scala lineare, uniti da linee per far risaltare la curva
    #plt.plot(degrees_er, probs_er, marker='s', linestyle='-', color=color_er, linewidth=2, alpha=0.9, label='Erdős-Rényi (Campana di Poisson)')
    plt.plot(degrees_ln, probs_ln, marker='o', linestyle='-', color=color_ln, linewidth=2, alpha=0.8, label='Lightning Network')

    plt.xlabel('Grado del Nodo $k$')

    plt.ylabel('Probabilità $P(k)$')
    
    # TRUCCO FONDAMENTALE: Zoomiamo sull'asse X per tagliare i Super Hub 
    # e concentrarci sulla parte centrale dove avviene la "normale" distribuzione
    plt.xlim(0, 25) 
    
    plt.grid(True, linestyle=':', alpha=0.7)
    plt.legend()
    plt.tight_layout()
    plt.savefig('8_Distribuzione_Gradi_Lineare.pdf')
    print("Salvato: 8_Distribuzione_Gradi_Lineare.pdf")

else:
    print(f"ERRORE: Assicurati che i file '{file_ln}' e '{file_er}' siano presenti nella cartella.")

print("Generazione completata!")