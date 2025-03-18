# IA04 - Systèmes Multi-Agents TD1-Vote en Go

Ce TD implémente les trois fonctions de l'agent client : créer un vote, voter et interroger les résultats du vote. Parallèlement, différentes règles de décompte des voix peuvent être définies : **« majorité », « borda », « approval », « condorcet », « copeland », « stv »**

5 exécutables (indépendants) sont fournis :

- `launch-all` permet de lancer plusieurs agents clients de tous types
- `launch-server` permet de lancer un agent de type serveur
- `lanch-client-newballot` permet de lancer un agent de type client pour créer un Bulletin
- `launch-client-vote` permet de lancer un agent de type client pour soumettre un Vote
- `launch-client-result` permet de lancer un agent de type client pour récupérer le Résultat

## Installation

Pour installer le paquet multiagentvote, utilisez la commande go get suivante :

```
go get gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote
```

Pour installer directement les exécutables :

```
go install  gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/cmd/launch-all
go install  gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/cmd/launch-server
go install  gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/cmd/launch-client-newballot
go install  gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/cmd/launch-client-vote
go install  gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/cmd/launch-client-result
```

## Usage

Avant d'exécuter l'agent client, veuillez exécuter **`launch-server`** pour activer l'agent serveur.

### Créer un Bulletin

Pour créer un nouveau bulletin de vote, suivez l'exemple ci-dessous pour modifier les paramètres et exécuter **`launch-client-newballot`**.

```
type RestClientAgent_NewBallot struct {
	id        string    // Identifiant ou unique représentant l'agent
	url       string    // URL du service REST que l'agent va consommer
	rule      string    // Règle de vote à appliquer
	deadline  time.Time // Date limite pour le vote
	voter_ids []string  // Identifiants des électeurs autorisés à participer au vote
	num_alts  int       // Nombre de candidats en lice
	tie_break []int     // Méthode pour départager en cas d'égalité
}
```

```
    ddl := "2023-10-30T23:05:08+04:00"
	t, _ := time.Parse(time.RFC3339, ddl)
	ag := client.NewRestClientAgent_NewBallot("id1", "http://localhost:8080", "borda", t, []string{"ag_id1", "ag_id2", "ag_id3"}, 2, []int{2, 1})
```

En même temps, il peut être créé directement dans REST.
Name | Value
-|-
rule|borda
deadline|2023-10-30T23:05:08+04:00
voter-ids|["ag_id1", "ag_id2", "ag_id3"]
#alts|2
tie-break|[2,1]

Parmi eux, le tie-break est facultatif.

### Soumettre un Vote

Pour soumettre un Vote, suivez l'exemple ci-dessous pour modifier les paramètres et exécuter **`launch-client-vote`**.

Si le temps de vote est expiré ou si l'électeur a déjà soumis un vote, ce comportement de vote échouera.

```
type RestClientAgent_Vote struct {
	id        string // Identifiant ou unique de l'agent
	url       string // URL du service REST que l'agent va utiliser
	agend_id  string // Nom de l'électeur
	ballot_id string // Identifiant du bulletin de vote
	prefs     []int  // Préférences de vote
	options   []int  // Options de vote
}
```

```
ag := client.NewRestClientAgent_Vote("id3", "http://localhost:8080", "ag_id1", "scrutin1", []int{4, 2, 3, 1}, []int{1})
```

Si sur REST, entrez comme suit:
Name | Value
-|-
agent-id|ag_id1
ballot-id|scrutin1
prefs|[4, 2, 3, 1]
options|[1]

Remarque :options est facultatif et permet de passer des renseignements supplémentaires (par exemple le seuil d'acceptation en approval)

### Récupérer le Résultat

Pour récupérer le Résultat, suivez l'exemple ci-dessous pour modifier les paramètres et exécuter **`launch-client-result`**.

Si le vote n'est pas encore clôturé, la requête échouera.

```
type RestClientAgent_Result struct {
	id        string //Identifiant ou unique de l'agent
	url       string //URL du service REST que l'agent va interroger
	ballot_id string // Identifiant du bulletin de vote
}
```

```
ag := client.NewRestClientAgent_Result("id2", "http://localhost:8080", "scrutin1")
```

Si sur REST, entrez comme suit:
Name | Value
-|-
ballot-id|scrutin1

### Commentaire

Depuis l'utilisation de SWFFactory et SCFFactory nous devons passer une fonction au format suivant：

```
func swf(p Profile) (Count, error)
func scf(p Profile) ([]Alternative, error)
```

Nous avons modifié les formats d'ApprovalSWF et ApprovalSCF comme suit pour répondre aux besoins d'utilisation de SWFFactory et SCFFactory.

```
func ApprovalSWF(thresholds []int) func(Profile) (Count, error)
func ApprovalSCF(thresholds []int) func(Profile) ([]Alternative, error)
```
