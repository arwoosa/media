## API Specifications

### REST API

#### Images Management
- **POST /api/v1/images**
  - **Description**: Upload images to Cloudflare
  - **Request**:
    - Content-Type: multipart/form-data
    - Parameters:
      - `images`: array of binary files (PNG, GIF, JPEG, WebP, SVG)
      - `metadata`: object containing tags and description
  - **Response**: Array of ImageResponse objects

- **GET /api/v1/images/{id}**
  - **Description**: Get image details
  - **Parameters**:
    - `id`: Image ID
  - **Response**: ImageResponse object

- **PUT /api/v1/images/{id}/permissions**
  - **Description**: Update image permissions
  - **Parameters**:
    - `id`: Image ID
  - **Request**: ImagePermission object
  - **Response**: Success message

- **POST /api/v1/images/{id}/process**
  - **Description**: Process image using Cloudflare Image API
  - **Parameters**:
    - `id`: Image ID
  - **Request**: ImageProcessingRequest object
  - **Response**: Processed ImageResponse object

- **DELETE /api/v1/images/{id}**
  - **Description**: Delete image with system retention (7 days)
  - **Parameters**:
    - `id`: Image ID
  - **Response**: Success message

#### Statistics
- **GET /api/v1/stats**
  - **Description**: Get system statistics
  - **Response**: SystemStats object containing:
    - Total images
    - Total storage usage
    - Usage by format
    - Traffic statistics
    - Storage usage details

## gRPC Service

### MediaService
```protobuf
service MediaService {
    // Upload images to Cloudflare
    rpc UploadImages (UploadRequest) returns (UploadResponse);
    
    // Get image details
    rpc GetImage (ImageRequest) returns (ImageResponse);
    
    // Update image permissions
    rpc UpdatePermissions (PermissionRequest) returns (PermissionResponse);
    
    // Process image using Cloudflare Image API
    rpc ProcessImage (ProcessRequest) returns (ProcessResponse);
    
    // Delete image with retention
    rpc DeleteImage (DeleteRequest) returns (DeleteResponse);
    
    // Get system statistics
    rpc GetStats (StatsRequest) returns (StatsResponse);
}
```

### Data Models

```protobuf
syntax = "proto3";

import "google/protobuf/timestamp.proto";

// Image format enum
enum ImageFormat {
  PNG = 0;
  GIF = 1;
  JPEG = 2;
  WEBP = 3;
  SVG = 4;
}

// CDN provider enum
enum CDNProvider {
  CLOUDFLARE = 0;
  LOCAL = 1;
}

// Image metadata
message Metadata {
  repeated string tags = 1;
  string description = 2;
}

// Access log entry
message AccessLog {
  string user_id = 1;
  google.protobuf.Timestamp accessed_at = 2;
  string ip_address = 3;
  string user_agent = 4;
}

// Image permissions
message Permissions {
  bool public = 1;
  repeated string allowed_users = 2;
  repeated string allowed_groups = 3;
  google.protobuf.Timestamp expires_at = 4;
  repeated AccessLog access_logs = 5;
}

// CDN status
message CDNStatus {
  bool cached = 1;
  string cdn_url = 2;
  int32 cache_ttl = 3;
  CDNProvider cdn_provider = 4;
}

// Traffic statistics
message TrafficStatistics {
  message DailyStats {
    string date = 1;
    int64 requests = 2;
    int64 bandwidth = 3;
  }
  repeated DailyStats daily = 1;
}

// Storage usage
message StorageUsage {
  int64 used = 1;
  int64 total = 2;
  double percentage = 3;
}

// System statistics
message SystemStats {
  int64 total_images = 1;
  int64 total_size = 2;
  map<string, int64> usage_by_format = 3;
  TrafficStatistics traffic_stats = 4;
  StorageUsage storage_usage = 5;
}

// Image response
message ImageResponse {
  string id = 1;
  string url = 2;
  ImageFormat format = 3;
  int64 size = 4;
  google.protobuf.Timestamp uploaded_at = 5;
  Metadata metadata = 6;
  Permissions permissions = 7;
  CDNStatus cdn_status = 8;
}
```

#### Enum Definitions
```protobuf
enum ImageFormat {
  PNG = 0;
  GIF = 1;
  JPEG = 2;
  WEBP = 3;
  SVG = 4;
}

enum CDNProvider {
  CLOUDFLARE = 0;
  LOCAL = 1;
}
```

## MongoDB Schema Design

### Collections

### 1. Images Collection
```javascript
{
    _id: ObjectId,  // MongoDB's auto-generated ID
    image_id: String,  // Unique identifier for the image
    cdn_url: String,      // Cloudflare Images URL
    format: String,   // enum: ['png', 'gif', 'jpeg', 'webp', 'svg']
    size: Number,     // File size in bytes
    uploaded_at: Date,
    metadata: {
        tags: [String],
        description: String,
        uploader_id: String,  // Reference to Users collection
        upload_ip: String
    },
    cloudflare_images: {
        account_id: String,     // Cloudflare account ID
        image_id: String,       // Cloudflare Images ID
        upload_status: String,  // enum: ['pending', 'success', 'failed']
        upload_error: String,   // Only present if upload failed
        upload_at: Date,
        last_modified: Date,
        expires_at: Date,
        analytics: {
            views: Number,         // Total views
            unique_views: Number,  // Unique views
            last_viewed: Date,     // Last viewed timestamp
            geographic: {
                countries: Map<String, Number>,  // Country views
                regions: Map<String, Number>,   // Region views
                cities: Map<String, Number>    // City views
            },
            devices: {
                desktop: Number,
                mobile: Number,
                tablet: Number
            },
            browsers: Map<String, Number>     // Browser distribution
        },
        security: {
            is_private: Boolean,   // Private image flag
            access_control: {
                ip_whitelist: [String],  // Allowed IP addresses
                auth_required: Boolean,  // Authentication required
                auth_method: String,     // enum: ['token', 'jwt', 'basic']
                auth_expiry: Date        // Authentication expiry
            },
            watermark: {
                enabled: Boolean,
                position: String,        // enum: ['top-left', 'top-right', 'bottom-left', 'bottom-right']
                text: String,
                color: String            // Hex color
            }
        },
        metadata: {
            title: String,
            description: String,
            alt_text: String,
            copyright: String,
            license: String,
            custom: Map<String, String>  // Custom metadata
        }
    },
    deleted_at: Date,  // Null if not deleted
    retention_until: Date,  // When the image will be permanently deleted
    version: Number    // Version control for image updates
}
```

### 2. Processing History Collection
```javascript
{
    _id: ObjectId,
    image_id: String,  // Reference to Images collection
    operation: String,  // Type of processing
    parameters: Object, // Processing parameters
    processed_at: Date,
    operator_id: String, // Who processed the image
    status: String,     // enum: ['pending', 'success', 'failed']
    error_message: String,  // Only present if status is 'failed'
    duration_ms: Number,    // Processing duration in milliseconds
    created_at: Date
}
```

### 3. Access Logs Collection
```javascript
{
    _id: ObjectId,
    image_id: String,  // Reference to Images collection
    user_id: String,   // Reference to Auth Service
    accessed_at: Date,
    ip_address: String,
    user_agent: String,
    request_method: String,  // HTTP method used
    response_status: Number, // HTTP status code
    duration_ms: Number,    // Request duration in milliseconds
    created_at: Date
}
```

### 4. Traffic Statistics Collection
```javascript
{
    _id: ObjectId,
    image_id: String,  // Reference to Images collection
    date: Date,        // Aggregated by day
    metrics: {
        total_requests: Number,
        successful_requests: Number,
        failed_requests: Number,
        average_response_time: Number,  // in milliseconds
        max_response_time: Number,
        min_response_time: Number,
        bandwidth_usage: Number,        // in bytes
        unique_users: Number,
        geo_distribution: {
            country: String,
            city: String,
            count: Number
        },
        device_distribution: {
            type: String,  // 'desktop'/'mobile'/'tablet'
            count: Number
        },
        browser_distribution: {
            name: String,
            version: String,
            count: Number
        }
    },
    created_at: Date,
    updated_at: Date
}
```

### 服務依賴

#### 主動依賴（Media Service 發起請求）
1. **認證服務（Auth Service）**
   - 通訊方式：JWT Token
   - 主要用途：
     - 使用者認證
     - 權限驗證
     - 請求簽名驗證

2. **Cloudflare**
   - 通訊方式：REST API
   - 主要用途：
     - 圖片存儲
     - 圖片處理
     - CDN 加速
     - 圖片優化

3. **通知服務（Notification Service）**
   - 通訊方式：REST API
   - 主要用途：
     - 圖片操作通知
     - 圖片審核通知
     - 系統通知

#### 被動依賴（其他服務發起請求）
1. **活動服務（Event Service）**
   - 通訊方式：gRPC
   - 主要用途：
     - 活動圖片查詢
     - 活動圖片統計
     - 活動圖片狀態

2. **社群服務（Community Service）**
   - 通訊方式：gRPC
   - 主要用途：
     - 社群圖片查詢
     - 社群圖片統計
     - 社群圖片狀態

3. **檢舉服務（Report Service）**
   - 通訊方式：REST API
   - 主要用途：
     - 圖片檢舉查詢
     - 圖片審核狀態
     - 圖片檢舉統計

### 依賴策略

#### 主動依賴策略
1. **請求管理**
   - 請求超時處理
   - 錯誤重試機制
   - 服務備份策略

2. **資料一致性**
   - 請求確認機制
   - 數據備份策略
   - 數據恢復機制

#### 被動依賴策略
1. **請求處理**
   - 請求驗證機制
   - 負載均衡
   - 錯誤回應處理

2. **效能優化**
   - 請求緩存
   - 數據預取
   - 數據壓縮

3. **安全控制**
   - 請求簽名驗證
   - 權限控制
   - 請求頻率限制

#### 服務隔離
1. **服務隔離**
   - 各服務獨立部署
   - 服務間通訊異步化
   - 服務狀態獨立監控

2. **錯誤處理**
   - 請求超時處理
   - 服務不可用備份
   - 錯誤重試機制

3. **資料一致性**
   - 服務間資料同步
   - 數據備份策略
   - 數據恢復機制
```

### 3. Ory Keto Integration

#### Permission Model
We'll use Ory Keto's permission model to manage image permissions:

1. **Namespace**: `images`
2. **Objects**: `image:<id>` (e.g., `image:123`)
3. **Relations**:
   - `view`: View image permission
   - `edit`: Edit image permission
   - `delete`: Delete image permission
   - `process`: Process image permission
   - `manage`: Manage image permissions

#### Example Permissions
```javascript
// User permission
{
    namespace: "images",
    object: "image:123",
    relation: "view",
    subject: {
        id: "user:456",
        type: "user"
    },
    expires_at: "2026-06-29T00:37:40+08:00"
}

// Group permission
{
    namespace: "images",
    object: "image:123",
    relation: "edit",
    subject: {
        id: "group:789",
        type: "group"
    }
}
```

#### Integration Flow

1. **Permission Checks**
```javascript
// Check if user can view image
GET /engines/ac/ory/roles/v0/permissions/check
{
    "subject": {
        "id": "user:456",
        "type": "user"
    },
    "namespace": "images",
    "object": "image:123",
    "relation": "view"
}
```

2. **Permission Updates**
```javascript
// Add permission
PUT /engines/ac/ory/roles/v0/permissions
{
    "subject": {
        "id": "user:456",
        "type": "user"
    },
    "namespace": "images",
    "object": "image:123",
    "relation": "view"
}

// Remove permission
DELETE /engines/ac/ory/roles/v0/permissions
{
    "subject": {
        "id": "user:456",
        "type": "user"
    },
    "namespace": "images",
    "object": "image:123",
    "relation": "view"
}
```

#### Best Practices

1. **Cache Strategy**
   - Implement local caching of permission checks
   - Use cache invalidation when permissions change
   - Set appropriate cache TTL based on permission type

2. **Error Handling**
   - Implement retry logic for permission checks
   - Handle Keto service downtime gracefully
   - Provide fallback mechanisms for critical operations

3. **Monitoring**
   - Monitor Keto service health
   - Track permission check latency
   - Log permission changes for audit
```

### 4. Statistics Collection
```javascript
{
    _id: ObjectId,
    type: String,      // enum: ['daily', 'weekly', 'monthly']
    date: Date,
    metrics: {
        total_images: Number,
        total_size: Number,
        usage_by_format: {
            png: Number,
            gif: Number,
            jpeg: Number,
            webp: Number,
            svg: Number
        },
        traffic: {
            requests: Number,
            bandwidth: Number,
            peak_hour: {
                hour: Number,
                requests: Number,
                bandwidth: Number
            }
        },
        storage: {
            used: Number,
            total: Number,
            percentage: Number
        },
        user_activity: {
            active_users: Number,
            uploads: Number,
            processes: Number
        }
    }
}
```

### Indexes

```javascript
// Images Collection
{
    // Primary indexes
    { image_id: 1 }: { unique: true },
    { url: 1 }: { unique: true },
    
    // Secondary indexes
    { uploaded_at: -1 },
    { metadata.uploader_id: 1 },
    { permissions.public: 1 },
    { permissions.expires_at: 1 },
    { deleted_at: 1 },
    { retention_until: 1 },
    { cdn_status.cached: 1 },
    { cdn_status.cdn_provider: 1 }
}

// Users Collection
{
    { username: 1 }: { unique: true },
    { email: 1 }: { unique: true },
    { role: 1 },
    { status: 1 },
    { last_login: -1 }
}

// Groups Collection
{
    { group_id: 1 }: { unique: true },
    { name: 1 },
    { status: 1 },
    { owner_id: 1 }
}

// Statistics Collection
{
    { date: 1 },
    { type: 1 },
    { metrics.storage.used: 1 }
}
```

### Data Relationships

1. **One-to-Many**:
   - Users can have many images (uploader_id in Images)
   - Groups can have many users (members array in Groups)
   
2. **Many-to-Many**:
   - Users can belong to multiple groups
   - Groups can have multiple users

### Schema Design Considerations

1. **Performance**:
   - Embedded access_logs and processing_history for better read performance
   - Proper indexing for common queries
   - Partitioning by date for statistics data

2. **Scalability**:
   - Sharding key consideration for Images collection
   - TTL indexes for temporary data
   - Proper use of arrays for permissions and groups

3. **Data Integrity**:
   - Referential integrity through IDs
   - Version control for image updates
   - Audit trails through access_logs

4. **Security**:
   - Permissions embedded in documents
   - Access control through user and group roles
   - Audit logging for all operations