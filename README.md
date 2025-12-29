# MinIO Alternative with example Golang

This repository demonstrates an alternative MinIO image setup and a simple Go client integration, intended for local development, experimentation, and as a foundation for further system integration.

## 📖 Related Article

This repository is a companion project for the following article:

👉 DevOps DIY [MinIO](https://andriantriputra.medium.com/devops-diy-minio-exploring-an-alternative-image-setup-and-go-client-integration-2e9280da397f): Exploring an Alternative Image Setup and Go Client Integration  


## Configuration
- Adjust ports and credentials directly in at compose
- Update environment variables as needed for your local or target environment

## ⚠️ Security Notes
Although MinIO can be started quickly using default credentials such as admin / admin@123, this setup is not safe for public or production environments.

Before deploying MinIO to a production server, consider the following:
1. Change Default Credentials
```bash
Update MINIO_ROOT_USER and MINIO_ROOT_PASSWORD in docker-compose.yml
to strong, unique values.
```

2. Enable HTTPS
```bash
Use TLS to encrypt traffic on ports 9000 and 9001.
MinIO supports both self-signed certificates and custom certificates
(e.g. from Let's Encrypt).
```

3. Restrict Network Access
```bash
If MinIO is only used internally, run it inside a private network
(e.g. Docker network or behind a reverse proxy).
```

4. Isolate Sensitive Storage
```bash
Use a dedicated volume or directory with restricted permissions
for the /data path.
```

5. Monitor Logs and Image Versions
```bash
Since community images may not receive automatic updates,
regularly rebuild from source and check MinIO security advisories.
```

These steps are simple but critical if you plan to use this setup in production or as part of a larger system, such as AI agent integration with n8n.

## Deployment
Start MinIO using Make:
```bash
$ make up
```

Check running containers:
```bash
$ make ps
```

View logs:
```bash
$ make logs
```

## Access
- Web Console → http://localhost:9001
- API → http://localhost:9000


## Reference
- [Hacker News discussion](https://news.ycombinator.com/item?id=45665452)
- [Reddit thread](https://www.reddit.com/r/selfhosted/comments/1oeibjk/minio_docker_image_with_the_classic_admin_web_ui/)
- [Related GitHub repository](https://github.com/Harsh-2002/MinIO)

## Author

Andrian Tri Putra
- [Medium](https://andriantriputra.medium.com/)
- [andriantp](https://github.com/andriantp)
- [AndrianTriPutra](https://github.com/AndrianTriPutra)

---

## License
Licensed under the Apache License 2.0