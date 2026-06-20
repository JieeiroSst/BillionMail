<div align="center">
  <a name="readme-top"></a>
  <h1><a href="https://www.billionmail.com/" target="_blank">BillionMail 📧</a></h1>

  <p><b>Mail Server tự host + nền tảng Email Marketing mã nguồn mở</b></p>

  [![][license-shield]][license-link] [![][docs-shield]][docs-link] [![][github-release-shield]][github-release-link] [![][github-stars-shield]][github-stars-link]

  Tài liệu dịch: [English](https://github.com/aaPanel/BillionMail) | [简体中文](README-zh_CN.md) | [日本語](README-ja.md) | [Türkçe](README-tr.md)
</div>

> Tài liệu này tập trung vào **cách triển khai và phát triển trên source code** của BillionMail. Giới thiệu sản phẩm/marketing xem tại [billionmail.com](https://www.billionmail.com/).

## Mục lục

1. [Tổng quan kiến trúc](#1-tổng-quan-kiến-trúc)
2. [Yêu cầu hệ thống](#2-yêu-cầu-hệ-thống)
3. [Cài đặt](#3-cài-đặt)
4. [Cấu hình (.env)](#4-cấu-hình-env)
5. [Cấu trúc source code](#5-cấu-trúc-source-code)
6. [Script quản lý `bm`](#6-script-quản-lý-bm)
7. [Quy trình phát triển](#7-quy-trình-phát-triển)
8. [Testing](#8-testing)
9. [Reverse proxy](#9-reverse-proxy)
10. [Quy ước khi sửa code](#10-quy-ước-khi-sửa-code)
11. [License](#11-license)

---

## 1. Tổng quan kiến trúc

BillionMail gồm 2 phần ghép lại qua Docker Compose ([docker-compose.yml](docker-compose.yml)):

**Mail stack** (xử lý gửi/nhận mail):
| Service | Image | Vai trò |
|---|---|---|
| `postfix-billionmail` | `billionmail/postfix` | MTA — gửi/nhận SMTP (25/465/587) |
| `dovecot-billionmail` | `billionmail/dovecot` | IMAP/POP3 (143/993/110/995), lưu mailbox |
| `rspamd-billionmail` | `billionmail/rspamd` | Lọc spam, DKIM signing |
| `webmail-billionmail` | `roundcube/roundcubemail` | Webmail UI, mount tại `/roundcube/` |

**Application stack** (quản trị, campaign, contact, tracking...):
| Service | Image | Vai trò |
|---|---|---|
| `core-billionmail` | `billionmail/core` | Backend Go (GoFrame) + frontend Vue đã build, expose cổng 80/443 |
| `pgsql-billionmail` | `postgres:17.4-alpine` | Database chính (mailbox, campaign, contact, log...) |
| `redis-billionmail` | `redis:7.4.2-alpine` | Cache, queue, rate-limit, session rspamd |

`core-billionmail` là phần bạn sẽ chỉnh sửa nhiều nhất khi phát triển: nó chứa toàn bộ logic nghiệp vụ (campaign, contact, domain, warmup, RBAC...) và giao tiếp với mail stack qua DB dùng chung + Docker socket (`/var/run/docker.sock`) để đọc trạng thái các container khác.

## 2. Yêu cầu hệ thống

**Để chạy bằng Docker (đủ dùng cho hầu hết trường hợp):**
- Linux x86_64 hoặc aarch64 (64-bit). Script cài đặt từ chối các kiến trúc khác.
- Docker + Docker Compose plugin (script `install.sh` tự cài nếu chưa có).
- Quyền root khi chạy `install.sh`.

**Để code trực tiếp lên backend/frontend:**
- Go ≥ 1.22 (xem [core/go.mod](core/go.mod), hiện khai báo `go 1.24.0`).
- Node.js + `pnpm` cho [core/frontend](core/frontend).
- Python 3 + `pytest` nếu chạy bộ test e2e ở [e2e/](e2e).

## 3. Cài đặt

### Cách 1 — Script cài đặt tự động (khuyến nghị)

```shell
cd /opt && git clone https://github.com/aaPanel/BillionMail && cd BillionMail
bash install.sh
# hoặc chỉ định domain/timezone:
bash install.sh --domain mail.example.com --TZ Asia/Ho_Chi_Minh
```

`install.sh` sẽ: kiểm tra kiến trúc/hệ điều hành → cài Docker nếu thiếu → sinh ngẫu nhiên `ADMIN_USERNAME`, `ADMIN_PASSWORD`, `SafePath`, `DBPASS`, `REDISPASS` → ghi vào `.env` → `docker compose up -d` → mở firewall cho các cổng mail/HTTP(S) → in ra địa chỉ truy cập và tài khoản admin.

Sau khi cài xong, truy cập theo địa chỉ in ra cuối log cài đặt:
```
https://<ip-server>:<HTTPS_PORT>/<SafePath>
```
(`SafePath` là một "cổng vào an toàn" sinh ngẫu nhiên, thay cho việc lộ `/login` công khai. `HTTPS_PORT` mặc định là `443` — khi đó có thể bỏ qua phần `:<HTTPS_PORT>`. Cả `SafePath` và `HTTPS_PORT` lấy trong file `.env`.)

Nếu cài trên máy local (Docker chạy ngay trên máy bạn đang mở trình duyệt), mở trực tiếp: [https://localhost](https://localhost) hoặc [https://localhost:8443](https://localhost:8443) nếu đã đổi `HTTPS_PORT`, tương tự với `http://localhost:<HTTP_PORT>` nếu chạy HTTP (thêm `/<SafePath>` vào sau).

> Ví dụ với cấu hình mặc định trong `.env` của repo này (`SafePath=billion`, `HTTPS_PORT=443`): mở **[https://localhost:443/billion](https://localhost:443/billion)**.

### Cách 2 — Docker Compose thủ công

```shell
cd /opt && git clone https://github.com/aaPanel/BillionMail && cd BillionMail
cp env_init .env
# chỉnh .env theo mục 4 trước khi chạy
docker compose up -d || docker-compose up -d
```

Dùng cách này khi bạn muốn tự kiểm soát từng biến môi trường, hoặc deploy trên hạ tầng có sẵn Docker/Compose.

### Cách 3 — One-click trên aaPanel

aaPanel → Docker → One-click install. Xem https://www.aapanel.com/new/download.html

## 4. Cấu hình (`.env`)

File `.env` (sinh từ [env_init](env_init)) là nguồn cấu hình duy nhất, được mount vào hầu hết container. Các biến quan trọng:

| Biến | Ý nghĩa |
|---|---|
| `ADMIN_USERNAME` / `ADMIN_PASSWORD` | Tài khoản đăng nhập quản trị |
| `SafePath` | Đường dẫn bí mật để truy cập trang quản trị (`https://host/<SafePath>`) |
| `BILLIONMAIL_HOSTNAME` | Hostname mail server (Postfix `myhostname`) — **luôn dùng dạng `mail.example.com`**, xem [mục 10](#10-quy-ước-khi-sửa-code) |
| `DBNAME` / `DBUSER` / `DBPASS` | PostgreSQL |
| `REDISPASS` | Redis |
| `SMTP_PORT`, `SMTPS_PORT`, `SUBMISSION_PORT`, `IMAP_PORT`, `IMAPS_PORT`, `POP_PORT`, `POPS_PORT` | Cổng mail |
| `HTTP_PORT` / `HTTPS_PORT` | Cổng truy cập trang quản trị |
| `WEB_BASE_PATH` | Base path khi nhúng sau reverse proxy theo path (ví dụ `/billionmail` khi chạy sau aaPanel) |
| `AAPANEL_SSO_SECRET` | Secret để aaPanel xin token SSO tin cậy — để trống nếu không dùng SSO |
| `TZ` | Timezone container |
| `IPV4_NETWORK` | Subnet nội bộ cho network Docker `billionmail-network` |
| `FAIL2BAN_INIT` | `y`/`n` — bật fail2ban chặn brute-force |
| `IP_WHITELIST_ENABLE` | `true`/`false` — chỉ cho IP whitelist truy cập console |
| `RETENTION_DAYS` | Số ngày giữ log/backup |

Sau khi sửa `.env`, áp dụng bằng `bm restart` hoặc `docker compose up -d` để container nhận giá trị mới.

> Một số tính năng gọi API bên ngoài qua biến môi trường runtime riêng (không thuộc `.env`/`env_init`, set trực tiếp cho process/container khi cần dùng tính năng đó):
> - `MOBILE_NOTIFY_URL` / `MOBILE_NOTIFY_TOKEN` — URL và token (tùy chọn, gửi kèm header `Authorization: Bearer <token>`) của backend mobile nhận push notification, dùng bởi `POST /api/notification/send_mobile` (xem [core/internal/service/notification](core/internal/service/notification)). Field `body` trong request là JSON array byte (ví dụ `[72,101,108,108,111]`), không phải base64 string — implement tại type `ByteArray` trong [core/api/notification/v1/notification.go](core/api/notification/v1/notification.go).

## 5. Cấu trúc source code

```
core/
├── internal/
│   ├── cmd/            # Entry point CLI, khởi tạo HTTP server, route, middleware
│   ├── controller/     # HTTP handler — 1 thư mục / domain (campaign, contact, domains, rbac...)
│   ├── service/        # Business logic — batch_mail, warmup, maillog_stat, rbac, fail2ban, acme...
│   ├── dao/            # Data access layer (sinh từ DB qua GoFrame CLI)
│   ├── model/entity/   # ORM entity
│   └── consts/         # Hằng số dùng chung
├── api/                # Định nghĩa route + struct request/response (input cho `gf gen ctrl`)
├── frontend/src/
│   ├── views/          # Trang (1 thư mục / feature)
│   ├── components/     # Component dùng lại
│   ├── store/          # Pinia store
│   ├── api/modules/    # API client gọi backend
│   ├── router/         # Vue Router + route theo module
│   ├── hooks/          # Composable
│   ├── utils/          # Tiện ích (base, data, time, storage)
│   ├── features/       # Component đặc thù theo feature (EmailEditor)
│   └── i18n/           # Đa ngôn ngữ (en, zh, ja)
├── template/           # Email template mặc định
└── manifest/           # Config GoFrame (config.yaml), Dockerfile build core
conf/                   # Config runtime cho mail stack (postfix, dovecot, rspamd, redis, webmail, fail2ban)
Dockerfiles/            # Dockerfile cho core, dovecot, postfix, rspamd
e2e/                    # Test end-to-end (pytest) — chạy lên toàn bộ stack thật
docs/                   # Tài liệu vận hành bổ sung (reverse proxy...)
install.sh / bm.sh / update.sh   # Cài đặt / quản trị / cập nhật
```

- Controller → `core/internal/controller/`, 1 thư mục/domain.
- Service → `core/internal/service/`, 1 thư mục/domain.
- API route → `core/api/`, 1 thư mục/domain.
- Frontend view → `core/frontend/src/views/`, 1 thư mục/feature.
- Test đặt cạnh file nguồn (`*_test.go`, `*.test.ts`).

## 6. Script quản lý `bm`

Sau khi cài bằng `install.sh`, lệnh `bm` (symlink tới `bm.sh`) có sẵn toàn hệ thống:

```shell
bm help                          # danh sách đầy đủ
bm default                       # xem lại tài khoản/URL đăng nhập mặc định
bm start | stop | restart | status | down | rebuild
bm top                           # xem tiến trình tất cả service
bm ps                            # xem container
bm service-top postfix           # tiến trình của 1 service cụ thể
bm log-file postfix              # log của 1 service (file log thật)
bm log-container postfix         # log của 1 container (docker logs)
bm restart-service postfix       # restart 1 service + container
bm change-port                   # đổi cổng truy cập quản trị
bm change-tz                     # đổi timezone
bm change-user / change-password
bm change-safe-path              # đổi lại "cổng vào an toàn"
bm change-apply-ssl-port         # đổi cổng dùng để xin SSL
bm cancel-ip-limit                # gỡ giới hạn truy cập theo IP (fail2ban)
bm clear                         # dọn image cũ
bm update                        # cập nhật BillionMail lên bản mới
```

## 7. Quy trình phát triển

### Backend (Go)

Backend dùng GoFrame v2; `Makefile` ở [core/Makefile](core/Makefile) include [core/hack/hack.mk](core/hack/hack.mk):

```shell
cd core
make dao        # sinh DAO/DO/Entity từ schema DB (gf gen dao)
make ctrl        # sinh controller từ định nghĩa trong api/ (gf gen ctrl)
make service     # sinh service interface (gf gen service)
make enums       # sinh enum từ code hiện có
```

Kiểm tra chất lượng code — **chạy sau mỗi lần sửa file**:
```shell
cd core && go vet ./... && gofmt -l .
```

Build & deploy bản dev của backend vào stack Docker đang chạy (không rebuild lại image):
```shell
# 1. Tạo container build (1 lần, theo go-build.sh):
docker run -d --name p-g-alpine --hostname p-g-alpine --restart=always \
  -v ./core:/opt/core \
  -v ./Dockerfiles/core/repositories:/etc/apk/repositories \
  alpine:3.20 tail -f /dev/null

# 2. Build + copy binary/manifest/public/template vào container core-billionmail đang chạy, rồi restart nó:
cd core && bash run_dev.sh
```
`run_dev.sh` tự nhận diện kiến trúc (amd64/arm64), build trong container `p-g-alpine`, copy file đã build sang container `billionmail-core-billionmail-1` và restart nó — dùng để lặp code nhanh mà không phải build lại Docker image.

### Frontend (Vue 3 + Rsbuild)

```shell
cd core/frontend
pnpm install
pnpm dev        # dev server, hot reload
```

Khi dev, request `/api/*` được proxy sang một backend thật theo cấu hình trong [core/frontend/build/config.ts](core/frontend/build/config.ts) (chọn bằng biến môi trường `SERVER_NAME`, mặc định `192.168.66.124`). **Trỏ nó sang backend đang chạy của bạn** (ví dụ instance Docker local hoặc server dev nội bộ) trước khi code, hoặc thêm entry mới vào `serverConfig`.

```shell
pnpm run lint        # eslint cho .vue/.ts/.tsx — chạy sau mỗi lần sửa file
pnpm run lint:fix
pnpm test            # vitest run
pnpm run test:watch
pnpm run build        # build production vào frontend/dist
```

## 8. Testing

```shell
# Go — short mode, không cần DB
cd core && go test -count=1 -short ./internal/service/...

# Frontend
cd core/frontend && pnpm test

# Cả hai
cd core && go test -count=1 -short ./internal/service/... && cd frontend && pnpm test
```

**E2E** ([e2e/](e2e), pytest, chạy lên stack thật + service giả lập):
```shell
docker compose -f docker-compose.yml -f docker-compose.test.yml up -d
cd e2e && pytest
```
`docker-compose.test.yml` thêm `mailpit` (bắt mail gửi đi để assert nội dung) và `mockserver` (wiremock, giả lập API ngoài) vào network `billionmail-network`.

## 9. Reverse proxy

Nếu đặt BillionMail sau Nginx/Caddy/Traefik (SSL termination, domain riêng...), xem cấu hình mẫu và header bắt buộc tại [docs/REVERSE_PROXY.md](docs/REVERSE_PROXY.md). Domain reverse proxy được lưu ở setting `reverse_proxy_domain` (bảng `bm_options`) — BillionMail dùng giá trị này để sinh tracking link (open/click) thay vì địa chỉ container nội bộ.

## 10. Quy ước khi sửa code

- **Không bao giờ dùng bare domain** cho hạ tầng mail (DNS, cert, DKIM, dedicated IP). Luôn lấy hostname qua `public.FormatMX(domain)` (ví dụ trả về `mail.example.com`).
- Khi thêm controller mới, phải thêm tên module vào danh sách `modules` trong RBAC tại [core/internal/service/middlewares/rbac.go](core/internal/service/middlewares/rbac.go), nếu không endpoint sẽ bị RBAC chặn.

## 11. License

BillionMail phát hành theo **AGPLv3** — xem [LICENSE](LICENSE). Báo lỗi/đề xuất tính năng: [GitHub Issues](https://github.com/aaPanel/BillionMail/issues). Báo lỗ hổng bảo mật: xem [SECURITY.md](SECURITY.md).

---

<!-- BillionMail official link -->
[docs-link]: https://www.billionmail.com/

<!-- BillionMail Other link-->
[license-link]: https://www.gnu.org/licenses/agpl-3.0.html
[github-release-link]: https://github.com/aaPanel/BillionMail/releases/latest
[github-stars-link]: https://github.com/aaPanel/BillionMail
[github-issues-link]: https://github.com/aaPanel/BillionMail/issues

<!-- Shield link-->
[docs-shield]: https://img.shields.io/badge/documentation-148F76
[github-release-shield]: https://img.shields.io/github/v/release/aaPanel/BillionMail
[github-stars-shield]: https://img.shields.io/github/stars/aaPanel/BillionMail?color=%231890FF&style=flat-square
[license-shield]: https://img.shields.io/github/license/aaPanel/BillionMail
