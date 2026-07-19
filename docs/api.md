# VPSM API Documentation

This document describes the REST API endpoints provided by the `vpsm-api` server.

## Managing the API Server

You can manage the background API server daemon using the `vpsm api` command suite:

```bash
# Start the background API server daemon
vpsm api start

# Check background daemon running status
vpsm api status

# View background server logs (support -f/--follow to tail in real-time)
vpsm api logs [-f]

# Gracefully restart the API daemon
vpsm api restart

# Stop the background API daemon
vpsm api stop

# Format, validate, and reload configs (automatically restarts running API daemons)
vpsm config reload
```

---

## Endpoint Reference

### 1. List All Servers
Retrieves a list of all servers registered in the database.

* **URL**: `/servers`
* **Method**: `GET`
* **Response Status**: `200 OK`
* **Response Body**:
```json
[
  {
    "id": 1,
    "uuid": "uuid-159.223.134.57",
    "name": "root@159.223.134.57",
    "host": "159.223.134.57",
    "port": 22,
    "username": "root",
    "auth_type": "password",
    "auth_secret": "deri@1234_#$",
    "provider": "DigitalOcean",
    "region": "nyc3",
    "os_family": "ubuntu",
    "os_version": "22.04",
    "cpu_model": "Intel(R) Xeon(R) Gold 6140 CPU @ 2.30GHz",
    "cpu_cores": 2,
    "ram_total": "4.0 GB",
    "disk_total": "80 GB",
    "created_at": "2026-07-16T14:04:49.000Z",
    "updated_at": "2026-07-16T14:04:49.000Z"
  },
  {
    "id": 2,
    "uuid": "uuid-52.54.164.79",
    "name": "ubuntu@52.54.164.79",
    "host": "52.54.164.79",
    "port": 22,
    "username": "ubuntu",
    "auth_type": "key",
    "auth_secret": "-----BEGIN OPENSSH PRIVATE KEY-----\n...",
    "provider": "AWS",
    "region": "us-east-1",
    "os_family": "ubuntu",
    "os_version": "20.04",
    "cpu_model": "Intel(R) Xeon(R) CPU E5-2676 v3 @ 2.40GHz",
    "cpu_cores": 1,
    "ram_total": "1.0 GB",
    "disk_total": "8 GB",
    "created_at": "2026-07-16T14:13:30.000Z",
    "updated_at": "2026-07-16T14:13:30.000Z"
  }
]
```

---

### 2. Scan Server Inventory
Triggers a background scanner job to collect installed software packages and system specs from the target host over SSH.

* **URL**: `/servers/:id/scan`
* **Method**: `POST`
* **Response Status**: `200 OK`
* **Response Body**:
```json
{
  "status": "scan initiated"
}
```

---

### 3. List Notes for a Server
Retrieves logs and notes associated with a particular server.

* **URL**: `/notes?server_id=:id`
* **Method**: `GET`
* **Response Status**: `200 OK`
* **Response Body**:
```json
[
  {
    "id": 1,
    "server_id": 1,
    "title": "System Reboot Info",
    "content": "Reboot performed successfully on Thursday, July 16, 2026.",
    "created_at": "2026-07-16T14:10:00Z",
    "updated_at": "2026-07-16T14:10:00Z"
  }
]
```

---

### 4. List Background Daemon Events
Retrieves chronological execution audits and event logs recorded by the background system scheduler and workers.

* **URL**: `/events`
* **Method**: `GET`
* **Response Status**: `200 OK`
```json
[]
```
