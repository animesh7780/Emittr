# Deployment Guide

## Quick Deployment Checklist

- [ ] All code committed to GitHub
- [ ] Environment variables configured
- [ ] Database migrations tested locally
- [ ] Docker images build successfully
- [ ] All tests passing
- [ ] README documentation complete
- [ ] Frontend production build tested

## Option 1: Railway Deployment (Recommended)

Railway offers free tier with PostgreSQL and easy deployment.

### 1. Prepare Repository
```bash
# Ensure all code is committed
git add .
git commit -m "Ready for deployment"
git push
```

### 2. Create Railway Account
- Go to https://railway.app
- Sign up with GitHub
- Connect your GitHub repository

### 3. Configure Services in Railway Dashboard

#### Database Service
1. Create PostgreSQL Database
2. Railway will provide `DATABASE_URL`
3. Copy to backend and analytics service

#### Environment Variables for Backend
```
DATABASE_URL=<from-railway>
KAFKA_BROKER=<kafka-service-from-railway>:9092
KAFKA_TOPIC=game_events
PORT=8080
ENVIRONMENT=production
```

#### Environment Variables for Analytics
```
DATABASE_URL=<from-railway>
KAFKA_BROKER=<kafka-service-from-railway>:9092
KAFKA_TOPIC=game_events
KAFKA_GROUP=analytics_group
```

#### Environment Variables for Frontend
```
VITE_API_URL=<backend-service-url>
```

### 4. Deploy
- Railway auto-deploys on push to main
- Monitor deployment in Railway dashboard
- Check logs for errors

### 5. Verify Deployment
```bash
# Test backend health
curl https://<your-railway-backend>/health

# Test leaderboard API
curl https://<your-railway-backend>/api/leaderboard

# Test frontend
Open https://<your-railway-frontend> in browser
```

---

## Option 2: Heroku Deployment

Heroku requires credit card (free tier is limited).

### 1. Install Heroku CLI
```bash
# macOS
brew tap heroku/brew && brew install heroku

# Windows
# Download from https://devcenter.heroku.com/articles/heroku-cli

# Verify
heroku --version
```

### 2. Login and Create Apps
```bash
heroku login

# Create backend app
heroku create your-app-name-backend
heroku buildpacks:add heroku/go --app=your-app-name-backend

# Create frontend app
heroku create your-app-name-frontend
heroku buildpacks:add heroku/nodejs --app=your-app-name-frontend
```

### 3. Add Databases

#### PostgreSQL
```bash
heroku addons:create heroku-postgresql:standard-0 --app=your-app-name-backend
# Heroku provides DATABASE_URL automatically
```

#### Kafka (CloudKarafka)
```bash
heroku addons:create cloudkarafka:giraffe --app=your-app-name-backend
# Heroku provides CLOUDKARAFKA_URL automatically
# Parse it: "KAFKA_BROKER=xxx.cloudkarafka.com:9092"
```

### 4. Set Environment Variables
```bash
# Backend
heroku config:set PORT=8080 --app=your-app-name-backend
heroku config:set ENVIRONMENT=production --app=your-app-name-backend
heroku config:set KAFKA_TOPIC=game_events --app=your-app-name-backend

# Frontend
heroku config:set VITE_API_URL=https://your-app-name-backend.herokuapp.com --app=your-app-name-frontend
```

### 5. Deploy

#### Backend
```bash
cd backend

# Create Procfile
echo "web: ./backend" > Procfile

# Deploy
heroku git:remote -a your-app-name-backend
git push heroku main
```

#### Frontend
```bash
cd frontend

# Create Procfile
echo "web: npm run build && npx serve -s dist -l \$PORT" > Procfile

# Deploy
heroku git:remote -a your-app-name-frontend
git push heroku main
```

### 6. Verify
```bash
heroku logs --tail --app=your-app-name-backend
heroku logs --tail --app=your-app-name-frontend
```

---

## Option 3: Docker Hub + Cloud Run

### 1. Build Docker Images
```bash
# Backend
cd backend
docker build -t your-username/4-in-a-row-backend:latest .
docker push your-username/4-in-a-row-backend:latest

# Frontend
cd frontend
docker build -t your-username/4-in-a-row-frontend:latest .
docker push your-username/4-in-a-row-frontend:latest

# Analytics
cd analytics
docker build -t your-username/4-in-a-row-analytics:latest .
docker push your-username/4-in-a-row-analytics:latest
```

### 2. Deploy to Google Cloud Run
```bash
# Deploy backend
gcloud run deploy 4-in-a-row-backend \
  --image your-username/4-in-a-row-backend:latest \
  --platform managed \
  --region us-central1 \
  --set-env-vars DATABASE_URL=<your-db-url>,KAFKA_BROKER=<your-kafka>

# Deploy frontend
gcloud run deploy 4-in-a-row-frontend \
  --image your-username/4-in-a-row-frontend:latest \
  --platform managed \
  --region us-central1
```

---

## Option 4: AWS Deployment

### Using ECS + RDS
```bash
# 1. Create RDS PostgreSQL instance
# 2. Create Elastic Container Service cluster
# 3. Create task definitions for backend, frontend, analytics
# 4. Create services from task definitions
# 5. Set up ALB (Application Load Balancer) for routing
# 6. Configure Route53 for custom domain
```

---

## Post-Deployment

### 1. Verify Services
```bash
# Frontend
Open https://your-frontend-url.com

# Backend
curl https://your-backend-url/health

# Leaderboard
curl https://your-backend-url/api/leaderboard
```

### 2. Monitor Performance
- Check error logs
- Monitor database connections
- Track Kafka consumer lag
- Monitor API response times

### 3. Setup CI/CD Pipeline

#### GitHub Actions Example
Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Backend
      run: |
        cd backend
        go build -o backend
    
    - name: Build Frontend
      run: |
        cd frontend
        npm install
        npm run build
    
    - name: Deploy to Railway
      env:
        RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
      run: |
        npm i -g @railway/cli
        railway up --service backend
        railway up --service frontend
        railway up --service analytics
```

---

## Monitoring & Logging

### Application Health Checks
```bash
# Backend health
curl https://your-backend-url/health

# Frontend loads
curl https://your-frontend-url

# Database connectivity
psql $DATABASE_URL -c "SELECT COUNT(*) FROM players;"
```

### Log Aggregation
- Railway: Built-in logs in dashboard
- Heroku: `heroku logs --tail`
- AWS: CloudWatch logs
- GCP: Cloud Logging

### Metrics to Track
- **API Response Time**: Should be < 200ms
- **Database Queries**: Monitor connection pool
- **Kafka Consumer Lag**: Should be < 5 seconds
- **Error Rate**: Should be < 0.1%
- **Active Connections**: Monitor WebSocket connections

---

## Troubleshooting Deployment

### Backend won't start
```
Check logs for database connection errors
Verify DATABASE_URL is set correctly
Ensure PostgreSQL is accessible from deployment region
```

### Frontend not connecting to backend
```
Check VITE_API_URL is correct
Verify CORS is enabled on backend
Check network policies in deployment platform
```

### Kafka connection fails
```
Verify KAFKA_BROKER is correct
Check firewall rules allow Kafka port
Ensure Kafka cluster is running
```

### Database migrations not running
```
Manual migration:
psql $DATABASE_URL < migrations.sql
Or restart backend service after schema changes
```

---

## Cost Estimation

### Railway (Recommended)
- Free tier: Very limited (for learning)
- PostgreSQL: $0.10/GB + compute
- Estimated: $10-20/month for basic usage

### Heroku
- PostgreSQL: $9-50+/month depending on size
- Dynos: $7-25+/month per app
- Kafka (CloudKarafka): $10+/month
- Estimated: $50-100+/month

### AWS
- RDS PostgreSQL: $15-100+/month
- EC2 instances: $5-50+/month per instance
- NAT Gateway: $32/month
- Estimated: $100+/month

### Google Cloud Run
- Pay per request (very cheap for low traffic)
- Cloud SQL: $10-100+/month
- Estimated: $15-50/month

---

## Scaling Considerations

### Database Scaling
- Add read replicas for heavy traffic
- Implement caching (Redis)
- Database connection pooling

### Backend Scaling
- Horizontal scaling with load balancer
- Session affinity for WebSocket connections
- Rate limiting

### Frontend Scaling
- CDN for static assets (CloudFlare)
- Image optimization
- Code splitting

### Kafka Scaling
- Increase partitions as throughput grows
- Consumer group scaling
- Retention policy management

---

## Security Checklist

- [ ] Environment variables are secrets (not in code)
- [ ] Database has strong password
- [ ] HTTPS/WSS enabled
- [ ] CORS properly configured
- [ ] SQL injection prevention (parameterized queries)
- [ ] Rate limiting enabled
- [ ] Input validation on all endpoints
- [ ] Sensitive data not logged
- [ ] Regular security updates
- [ ] Monitoring and alerting configured

---

**Need Help?**

Check the logs first:
```
Railway: Dashboard → Logs
Heroku: heroku logs --tail
```

Common issues are usually:
1. Environment variables not set
2. Database not accessible
3. Kafka broker incorrect
4. Port conflicts
5. CORS issues

Good luck with deployment! 🚀
