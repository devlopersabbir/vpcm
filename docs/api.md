# VPSM API Documentation

This document describes the REST API endpoints provided by the `vpsm-api` server.

## Starting the API Server

By default, the API server can be started on the host and port defined in your configuration (usually `127.0.0.1:8080` or via environment variables).

```bash
vpsm-api
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
    "provider": "",
    "os_family": "",
    "os_version": "",
    "tags": null,
    "software": null,
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
    "provider": "",
    "os_family": "",
    "os_version": "",
    "tags": null,
    "software": null,
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
