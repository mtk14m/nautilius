# Notilius IDP - Developer Platform Self-Service

> Plateforme développeur complète avec self-service portal, construite en mode SRE pour homelab et portfolio.

## 🎯 Objectif du Projet

Construire une plateforme développeur (Developer Platform) de niveau professionnel permettant aux développeurs de:
- Provisionner des ressources de manière autonome (namespaces, clusters, services)
- Gérer leurs applications via un portal self-service
- Monitorer leurs services en temps réel
- Suivre les coûts et quotas
- Accéder à la documentation et aux APIs

## 🏗️ Architecture

Voir [ARCHITECTURE.md](./ARCHITECTURE.md) pour les détails complets de l'architecture.

### Stack Technologique

- **Infrastructure**: k3s (Kubernetes)
- **GitOps**: ArgoCD
- **Backend API**: Go (API Gateway) + Python (Microservices)
- **Frontend**: React/Next.js
- **Observability**: Prometheus, Grafana, Loki, Tempo
- **Database**: PostgreSQL
- **Security**: OIDC (Keycloak/Dex), Vault

## 📚 Documentation

- **[ARCHITECTURE.md](./ARCHITECTURE.md)**: Architecture détaillée et composants
- **[ROADMAP.md](./ROADMAP.md)**: Roadmap complète avec phases d'apprentissage
- **[LEARNING_GUIDE.md](./LEARNING_GUIDE.md)**: Guide d'apprentissage pour devenir Platform Engineer SRE

## 🚀 Démarrage Rapide

### Prérequis

- ✅ VPN configuré
- ✅ k3s installé et fonctionnel
- `kubectl` configuré pour accéder au cluster
- `helm` installé
- `git` installé

### Vérification de l'Environnement

```bash
# Vérifier l'accès au cluster
kubectl cluster-info
kubectl get nodes

# Vérifier que k3s est bien installé
kubectl get pods -n kube-system
```

### Prochaines Étapes

1. **Phase 1: Fondations Infrastructure** (Semaine 1-2)
   - [ ] Structurer le projet
   - [ ] Configurer l'ingress (Traefik)
   - [ ] Installer ArgoCD
   - [ ] Configurer cert-manager pour TLS

2. **Phase 2: Observability Stack** (Semaine 2-3)
   - [ ] Installer Prometheus Operator
   - [ ] Configurer Grafana
   - [ ] Déployer Loki pour les logs
   - [ ] Configurer l'alerting

3. **Phase 3: API Backend** (Semaine 3-5)
   - [ ] Créer l'API Gateway en Go
   - [ ] Développer les services Python
   - [ ] Configurer PostgreSQL
   - [ ] Implémenter l'authentification OIDC

4. **Phase 4: Self-Service Portal** (Semaine 5-7)
   - [ ] Créer le frontend React/Next.js
   - [ ] Intégrer avec l'API
   - [ ] Implémenter le dashboard développeur
   - [ ] Ajouter la gestion de projets

5. **Phase 5: Advanced Features** (Semaine 7-10)
   - [ ] Développer le CLI tool en Go
   - [ ] Créer les SDKs (Python, Go)
   - [ ] Documenter complètement
   - [ ] Mettre en place CI/CD

6. **Phase 6: Production Ready** (Semaine 10-12)
   - [ ] Tests end-to-end
   - [ ] Security hardening
   - [ ] Performance optimization
   - [ ] Documentation finale

## 📖 Guide d'Apprentissage

Ce projet est conçu pour vous transformer en **Platform Engineer** et **SRE**. 

**Approche recommandée**:
1. Lire le [LEARNING_GUIDE.md](./LEARNING_GUIDE.md) pour comprendre les concepts
2. Suivre la [ROADMAP.md](./ROADMAP.md) phase par phase
3. Implémenter chaque composant en comprenant le "pourquoi"
4. Documenter vos apprentissages

## 🎓 Concepts Clés à Maîtriser

- **Kubernetes**: Pods, Services, Deployments, RBAC, Custom Resources
- **GitOps**: ArgoCD, Application CRD, Sync strategies
- **Observability**: Prometheus, Grafana, Loki, OpenTelemetry
- **Microservices**: API Gateway, Service Discovery, Communication patterns
- **Go**: Concurrency, Kubernetes client-go, gRPC
- **Python**: FastAPI, Async/await, Kubernetes client
- **SRE**: SLI/SLO/SLA, Error Budget, Monitoring, Incident Response

## 📁 Structure du Projet

```
notilius-idp/
├── infrastructure/     # Terraform, scripts d'infra
├── gitops/            # Configurations GitOps
│   ├── base/
│   ├── apps/
│   └── environments/
├── platform-api/      # API Gateway Go
├── services/          # Microservices Python
│   ├── provisioning/
│   ├── monitoring/
│   ├── billing/
│   └── notifications/
├── portal/           # Frontend React/Next.js
├── cli/              # CLI Tool Go
├── docs/             # Documentation
└── scripts/          # Scripts utilitaires
```

## 🔒 Sécurité

- Secrets gérés via External Secrets Operator ou Vault
- RBAC fine-grained pour tous les services
- Network Policies pour isolation réseau
- Scanning des images avec Trivy
- TLS partout (cert-manager)

## 📊 Métriques de Succès

- Temps de provisionnement d'un namespace: < 2 minutes
- Disponibilité de la plateforme: 99.9% (SLO)
- Latence API p95: < 200ms
- Taux d'erreur: < 0.1%

## 🤝 Contribution

Ce projet est pour votre apprentissage personnel et portfolio. 

**Bonnes pratiques**:
- Commits atomiques avec messages clairs
- Documentation de chaque décision importante
- Tests pour chaque composant
- Code propre et maintenable

## 📝 Notes d'Apprentissage

Créez un fichier `NOTES.md` pour documenter:
- Concepts appris
- Décisions d'architecture (ADRs)
- Problèmes rencontrés et solutions
- Ressources utiles

## 🔗 Ressources Utiles

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Google SRE Book](https://sre.google/books/)
- [CNCF Landscape](https://landscape.cncf.io/)

## 📅 Timeline

- **Semaine 1-2**: Fondations Infrastructure
- **Semaine 2-3**: Observability Stack
- **Semaine 3-5**: API Backend
- **Semaine 5-7**: Self-Service Portal
- **Semaine 7-10**: Advanced Features
- **Semaine 10-12**: Production Ready

---

**Bon apprentissage ! 🚀**

*Rappel: Prenez le temps de comprendre chaque concept avant de passer au suivant. La qualité prime sur la vitesse.*



