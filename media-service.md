# Media Service 系統需求說明

## 核心功能

### 圖片上傳功能
**標題**：圖片上傳功能

**User Story**：
As a registered user
I want to upload images directly to Cloudflare
so that I can store and share visual content efficiently.

**業務背景**：
使用者需要能夠直接將圖片上傳到 Cloudflare，以獲得更好的圖片處理和加速效果。

**驗收標準**：
- 支援多種圖片格式（PNG, GIF, JPEG, WebP, SVG）
- 支援批量上傳到 Cloudflare
- 支援拖放上傳到 Cloudflare
- 支援圖片預覽
- 支援圖片格式驗證
- 支援圖片大小限制
- 支援 Cloudflare 令牌驗證
- 支援圖片 URL 管理
- 支援圖片版本控制
- 支援圖片使用統計

**估算**：3 Story Points

### 圖片權限管理功能
**標題**：圖片權限管理功能

**User Story**：
As a registered user
I want to set image permissions
so that I can control who can view my images.

**業務背景**：
使用者需要能夠設定圖片的查看權限，以保護圖片內容的安全性。

**驗收標準**：
- 支援圖片權限設定
- 支援特定用戶權限設定
- 支援用戶群組權限設定
- 支援權限繼承
- 支援權限修改
- 支援權限查詢
- 支援權限日誌記錄
- 支援權限驗證
- 支援權限統計

**估算**：3 Story Points

### 圖片自動刪除功能
**標題**：圖片自動刪除功能

**User Story**：
As a registered user
I want to delete images with system retention
so that I can manage storage space efficiently.

**業務背景**：
使用者需要能夠刪除圖片，並由系統保留一週後自動刪除，以管理儲存空間。

**驗收標準**：
- 支援圖片刪除
- 支援系統保留一週
- 支援刪除日誌記錄
- 支援刪除通知
- 支援刪除統計
- 支援刪除恢復（保留期內）
- 支援刪除驗證
- 支援刪除監控
- 支援刪除策略配置

**估算**：2 Story Points

### 圖片處理功能
**標題**：圖片處理功能

**User Story**：
As a registered user
I want to process images
so that I can optimize visual content for display.

**業務背景**：
使用者需要能夠處理圖片，以適應不同顯示需求。

**驗收標準**：
- 支援圖片裁剪（Cloudflare Image API）
- 支援圖片旋轉
- 支援圖片調整（亮度、對比度、飽和度）
- 支援圖片過濾（模糊、銳化、黑白等）
- 支援圖片壓縮
- 支援圖片水印
- 支援圖片格式轉換
- 支援圖片元數據編輯
- 支援圖片品質調整
- 支援圖片尺寸調整

**估算**：3 Story Points

### CDN 管理功能
**標題**：CDN 管理功能

**User Story**：
As an administrator
I want to manage CDN settings
so that I can optimize content delivery.

**業務背景**：
為了提高內容傳輸效率，需要管理 CDN 設置。

**驗收標準**：
- 支援 CDN 配置
- 支援 CDN 緩存
- 支援 CDN 監控
- 支援 CDN 優化
- 支援 CDN 故障切換
- 支援 CDN 成本控制
- 支援 CDN 緩存策略
- 支援 CDN 回源策略

**估算**：3 Story Points

### 圖片流量統計功能
**標題**：圖片流量統計功能

**User Story**：
As a business user
I want to track image traffic statistics
so that I can optimize media usage and manage costs.

**業務背景**：
商業用戶需要了解圖片的使用情況和流量，以便優化媒體資源和控制成本。

**驗收標準**：
- 支援圖片訪問統計
- 支援圖片流量統計
- 支援圖片使用趨勢分析
- 支援圖片熱點分析
- 支援圖片地域分布統計
- 支援圖片時間段統計
- 支援圖片格式統計
- 支援圖片大小統計
- 支援圖片來源統計
- 支援圖片使用排名
- 支援統計報表導出
- 支援統計數據導入
- 支援統計數據備份
- 支援統計數據恢復
- 支援統計數據清理

**估算**：4 Story Points

### 數據管理功能
**標題**：數據管理功能

**User Story**：
As a system administrator
I want to manage data storage
so that I can ensure data safety and efficiency.

**業務背景**：
為了確保數據的安全和效率，需要有效的數據管理。

**驗收標準**：
- 支援多種儲存方式（本地/Cloudflare）
- 支援數據恢復
- 支援數據監控
- 支援數據統計
- 支援數據備份
- 支援數據清理
- 支援數據優化
- 支援數據配額管理
- 支援數據成本控制
- 支援數據使用統計

**估算**：3 Story Points

### 安全管理功能
**標題**：安全管理功能

**User Story**：
As a system administrator
I want to implement security measures
so that I can protect system data and privacy.

**業務背景**：
為了確保系統的安全性，需要實現完善的安全措施。

**驗收標準**：
- 支援數據加密（Cloudflare KSM）
- 支援存取控制
- 支援檔案驗證
- 支援安全備份
- 支援安全監控
- 支援安全審計
- 支援安全策略配置
- 支援安全事件處理
- 支援安全日誌記錄
- 支援安全風險評估

**估算**：3 Story Points

### 系統整合功能
**標題**：系統整合功能

**User Story**：
As a system developer
I want to implement service integration
so that I can support cross-service operations.

**業務背景**：
為了實現系統間的協作，需要實現服務整合。

**驗收標準**：
- 與 Event Service 整合（gRPC）
- 與 Community Service 整合（gRPC）
- 與 Report Service 整合（REST API）
- 與 Notification Service 整合（REST API）
- 支援服務發現
- 支援服務路由
- 支援服務監控
- 支援服務備份
- 支援服務恢復
- 支援服務狀態監控

**估算**：3 Story Points

## 權限說明

### 使用者權限
- 一般使用者（User）
  - 圖片上傳權限
  - 圖片查詢權限
  - 圖片收藏權限
  - 圖片刪除權限（保留期內可恢復）

### 商業用戶權限（BusinessUser）
  - 媒體管理權限
  - 圖片權限設定
  - 圖片使用統計
  - 圖片備份管理

### 系統管理員權限（Admin）
  - 系統設定權限
  - 安全管理權限
  - 統計查詢權限
  - CDN 管理權限
  - 媒體審核權限
  - 使用者權限管理

## 使用情境

### 圖片管理
- 圖片上傳
  - 使用者：選擇圖片
  - 系統：驗證格式
  - 系統：檢查大小
  - 系統：處理壓縮

- 圖片處理
  - 系統：裁切圖片
  - 系統：壓縮圖片
  - 系統：添加水印
  - 系統：保存處理

- 圖片存取
  - 使用者：查詢圖片
  - 使用者：下載圖片
  - 主辦方：刪除圖片
  - 系統：記錄操作



## 技術要求

### 技術框架
- 使用 Cobra 構建 CLI 工具
- 使用 Viper 管理配置
- 使用 Gin 構建 REST API
- 使用 gRPC 構建微服務通訊

### Cloudflare 集成
- Cloudflare Image API 集成
- 圖片轉碼與優化
- 圖片快取策略
- 圖片加速配置

### 服務依賴
#### 主動依賴（Media Service 發起請求）
1. **Auth Service**
   - 通訊方式：JWT Token
   - 主要用途：
     - 使用者身份驗證
     - 權限驗證
     - 請求簽名驗證

2. **Cloudflare**
   - 通訊方式：REST API
   - 主要用途：
     - 圖片存儲
     - 圖片處理
     - CDN 加速
     - 圖片優化

3. **Notification Service**
   - 通訊方式：REST API
   - 主要用途：
     - 圖片操作通知
     - 圖片審核通知
     - 系統通知

#### 被動依賴（其他服務發起請求）
1. **Event Service**
   - 通訊方式：gRPC
   - 主要用途：
     - 活動圖片查詢
     - 活動圖片統計
     - 活動圖片狀態

2. **Community Service**
   - 通訊方式：gRPC
   - 主要用途：
     - 社群圖片查詢
     - 社群圖片統計
     - 社群圖片狀態

3. **Report Service**
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

### 依賴策略
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

### MongoDB 集成
- MongoDB 數據庫整合
- MongoDB 資料模型設計
- MongoDB 索引優化
- MongoDB 數據備份
- MongoDB 數據恢復
- MongoDB 數據監控
- MongoDB 數據統計

### 數據管理
- 支援多種儲存方式（本地/Cloudflare）
- 支援數據恢復
- 支援數據監控
- 支援數據統計

### 安全要求
- 支援數據加密（Cloudflare KSM）
- 支援存取控制
- 支援檔案驗證
- 支援安全備份
- 支援安全監控

### 開發工具
- Apidog 生成 API 文件
- Protobuf 定義 gRPC 接口
- Go Modules 管理依賴
- Docker 容器化部署
- CI/CD 集成

### 系統整合
- 與 Event Service 整合（gRPC）
- 與 Community Service 整合（gRPC）
- 與 Report Service 整合（REST API）
- 與 Notification Service 整合（REST API）

### 開發工具
- Swagger 生成 API 文件
- Protobuf 定義 gRPC 接口
- Go Modules 管理依賴
- Docker 容器化部署
- CI/CD 集成

## 介面需求

### 媒體管理介面
- 支援媒體上傳
- 支援媒體預覽
- 支援媒體編輯
- 支援媒體管理

### 儲存管理介面
- 支援儲存配置
- 支援儲存監控
- 支援儲存備份
- 支援儲存恢復

### CDN 管理介面
- 支援 CDN 配置
- 支援 CDN 監控
- 支援 CDN 優化
- 支援 CDN 統計

### 統計報表介面
- 支援媒體統計
- 支援儲存統計
- 支援 CDN 統計
- 支援數據報告

---

> 如需設計 API 文件、資料表結構、畫面 Wireframe 或測試驗收項目，請聯繫系統分析師。
