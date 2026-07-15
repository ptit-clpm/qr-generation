**HỌC VIỆN CÔNG NGHỆ BƯU CHÍNH VIỄN THÔNG**

**BÁO CÁO ĐẶC TẢ YÊU CẦU PHẦN MỀM**

**SOFTWARE REQUIREMENTS SPECIFICATION (SRS)**

**ĐỀ TÀI: WEBSITE TẠO VÀ QUẢN LÝ MÃ QR-CODE**

**Tên dự án: QR Generator - QR Studio**

Phiên bản tài liệu: 1.0

Ngày lập: 05/07/2026

Nhóm thực hiện: ................................................

Giảng viên hướng dẫn: ...........................................

Môn học: Đảm bảo chất lượng phần mềm

# **LỊCH SỬ PHIÊN BẢN TÀI LIỆU**

| **Phiên bản** | **Ngày**   | **Người thực hiện** | **Nội dung thay đổi**                                   |
|---------------|------------|---------------------|---------------------------------------------------------|
| 1.0           | 05/07/2026 | Nhóm dự án          | Tạo tài liệu SRS ban đầu cho dự án website tạo QR-Code. |

# **MỤC LỤC**

1\. Giới thiệu

2\. Mô tả tổng quan hệ thống

3\. Đối tượng sử dụng và phân quyền

4\. Yêu cầu chức năng

5\. Đặc tả Use Case chính

6\. Yêu cầu dữ liệu và thiết kế cơ sở dữ liệu

7\. Yêu cầu giao diện

8\. Yêu cầu phi chức năng

9\. Quy tắc nghiệp vụ

10\. Ràng buộc, giả định và phụ thuộc

11\. Tiêu chí nghiệm thu

12\. Phụ lục

*Ghi chú: Có thể cập nhật lại mục lục tự động trong Microsoft Word nếu
bổ sung thêm nội dung hoặc thay đổi số trang.*

# **DANH MỤC TỪ VIẾT TẮT**

| **Từ viết tắt** | **Ý nghĩa**                         | **Diễn giải**                                   |
|-----------------|-------------------------------------|-------------------------------------------------|
| SRS             | Software Requirements Specification | Tài liệu đặc tả yêu cầu phần mềm.               |
| QR              | Quick Response                      | Mã phản hồi nhanh, dạng mã vạch hai chiều.      |
| URL             | Uniform Resource Locator            | Đường dẫn đến tài nguyên trên Internet.         |
| UI              | User Interface                      | Giao diện người dùng.                           |
| UX              | User Experience                     | Trải nghiệm người dùng.                         |
| API             | Application Programming Interface   | Giao diện lập trình ứng dụng.                   |
| CRUD            | Create, Read, Update, Delete        | Các thao tác tạo, đọc, cập nhật và xóa dữ liệu. |
| ERD             | Entity Relationship Diagram         | Mô hình thực thể liên kết.                      |
| MVP             | Minimum Viable Product              | Phiên bản sản phẩm tối thiểu có thể sử dụng.    |

# **1. GIỚI THIỆU**

## **1.1. Mục đích tài liệu**

Tài liệu này đặc tả yêu cầu phần mềm cho dự án Website tạo và quản lý mã
QR-Code, tên gợi ý là QR Generator - QR Studio. Tài liệu được sử dụng
làm cơ sở thống nhất giữa nhóm phát triển, người kiểm thử, giảng viên
hướng dẫn và các bên liên quan trong quá trình phân tích, thiết kế,
triển khai và nghiệm thu hệ thống.

Nội dung tài liệu tập trung mô tả phạm vi hệ thống, đối tượng sử dụng,
yêu cầu chức năng, yêu cầu phi chức năng, yêu cầu dữ liệu, phân quyền,
quy tắc nghiệp vụ và tiêu chí nghiệm thu.

## **1.2. Phạm vi sản phẩm**

Hệ thống là một website trực tuyến cho phép người dùng tạo mã QR theo
nhiều loại nội dung khác nhau, tùy chỉnh giao diện QR, tải xuống mã QR,
lưu mã QR vào tài khoản và theo dõi thống kê lượt quét đối với Dynamic
QR.

Hệ thống hỗ trợ mô hình hai gói sử dụng: Free và Pro. Gói Free phục vụ
nhu cầu cơ bản, còn gói Pro mở khóa các chức năng nâng cao như Dynamic
QR, chỉnh sửa nội dung sau khi tạo, thêm logo, tải nhiều định dạng và
xem thống kê lượt quét.

## **1.3. Mục tiêu của dự án**

- Xây dựng website tạo mã QR nhanh, dễ dùng và phù hợp với nhiều nhóm
  người dùng.

- Hỗ trợ các loại QR phổ biến như URL, Text, WiFi, vCard, Email, SMS,
  Location, Social, PDF và Menu.

- Cho phép tùy chỉnh thiết kế mã QR như màu sắc, logo, kiểu mắt QR, kiểu
  chấm, khung và template.

- Cho phép người dùng đăng ký, đăng nhập, quản lý danh sách QR đã tạo.

- Xây dựng phân quyền Free/Pro/Admin rõ ràng.

- Hỗ trợ Dynamic QR cho phép chỉnh sửa nội dung đích sau khi mã đã được
  tạo.

- Ghi nhận và thống kê lượt quét mã QR cho người dùng Pro.

- Tích hợp thanh toán nâng cấp gói Pro.

- Cung cấp trang quản trị cho Admin theo dõi người dùng, QR, gói dịch
  vụ, giao dịch và template.

## **1.4. Phạm vi áp dụng**

- Cá nhân: tạo QR thông tin liên hệ, mạng xã hội, văn bản, email hoặc vị
  trí.

- Cửa hàng, quán cà phê, nhà hàng: tạo QR menu, QR WiFi, QR chương trình
  khuyến mãi.

- Doanh nghiệp nhỏ và người làm marketing: tạo QR cho chiến dịch quảng
  cáo, poster, banner, tờ rơi và bao bì sản phẩm.

- Trường học và tổ chức giáo dục: tạo QR tài liệu học tập, điểm danh,
  khảo sát, sự kiện hoặc hội thảo.

## **1.5. Định nghĩa sản phẩm**

Website QR Generator - QR Studio là hệ thống hỗ trợ tạo QR-Code trực
tuyến, trong đó người dùng có thể tạo QR tĩnh hoặc QR động. QR tĩnh mã
hóa trực tiếp nội dung vào hình QR, còn QR động chứa một đường dẫn trung
gian của hệ thống để ghi nhận lượt quét và chuyển hướng đến nội dung
đích hiện tại.

# **2. MÔ TẢ TỔNG QUAN HỆ THỐNG**

## **2.1. Bối cảnh hệ thống**

Mã QR ngày càng được sử dụng nhiều trong đời sống và kinh doanh. Tuy
nhiên, nhiều công cụ tạo QR miễn phí chỉ hỗ trợ tạo QR đơn giản, không
có quản lý lịch sử, không có thống kê và không cho phép chỉnh sửa nội
dung sau khi đã in mã QR. Vì vậy, hệ thống được xây dựng để giải quyết
nhu cầu tạo, quản lý và đo lường hiệu quả sử dụng QR trong một nền tảng
thống nhất.

## **2.2. Mô hình tổng quan**

Hệ thống được thiết kế theo mô hình Client - Server gồm các thành phần
chính:

- Frontend: giao diện website cho khách truy cập, người dùng Free, người
  dùng Pro và Admin.

- Backend API: xử lý nghiệp vụ đăng nhập, tạo QR, quản lý QR, phân
  quyền, thanh toán, thống kê.

- Database: lưu thông tin người dùng, mã QR, gói dịch vụ, thanh toán,
  template và lượt quét.

- Payment Gateway: xử lý thanh toán nâng cấp gói Pro.

- QR Redirect Service: xử lý truy cập Dynamic QR, ghi nhận lượt quét và
  chuyển hướng.

## **2.3. Nhóm chức năng chính**

| **Nhóm chức năng**    | **Mô tả**                                                                                               |
|-----------------------|---------------------------------------------------------------------------------------------------------|
| Xác thực và tài khoản | Đăng ký, đăng nhập, đăng xuất, đổi mật khẩu, cập nhật hồ sơ cá nhân.                                    |
| Tạo QR                | Tạo QR theo nhiều loại nội dung như URL, Text, WiFi, vCard, Email, SMS, Location.                       |
| Tùy chỉnh QR          | Tùy chỉnh màu sắc, logo, kiểu mắt, kiểu chấm, khung, kích thước và template.                            |
| Quản lý QR            | Xem danh sách, tìm kiếm, lọc, xem chi tiết, xóa, tải xuống và sao chép QR.                              |
| Dynamic QR            | Tạo đường dẫn trung gian, chỉnh sửa nội dung đích và chuyển hướng khi quét QR.                          |
| Thống kê              | Ghi nhận và hiển thị lượt quét theo thời gian, thiết bị, trình duyệt, hệ điều hành và vị trí tương đối. |
| Gói dịch vụ           | Quản lý Free/Pro, giới hạn chức năng, nâng cấp và gia hạn.                                              |
| Thanh toán            | Tạo giao dịch, xử lý kết quả từ cổng thanh toán và kích hoạt Pro.                                       |
| Quản trị              | Quản lý người dùng, QR, gói, giao dịch, template, nội dung và log hệ thống.                             |

## **2.4. Môi trường vận hành**

- Trình duyệt web hiện đại như Chrome, Edge, Firefox, Safari.

- Thiết bị người dùng gồm máy tính, máy tính bảng và điện thoại.

- Backend triển khai trên server hoặc cloud có hỗ trợ API REST.

- Database sử dụng MySQL hoặc PostgreSQL.

- Có kết nối Internet để tạo, lưu, tải và chuyển hướng Dynamic QR.

## **2.5. Công nghệ đề xuất**

| **Thành phần** | **Công nghệ đề xuất**                       | **Ghi chú**                                        |
|----------------|---------------------------------------------|----------------------------------------------------|
| Frontend       | ReactJS hoặc Next.js                        | Xây dựng giao diện người dùng và dashboard.        |
| CSS/UI         | Tailwind CSS hoặc Bootstrap                 | Tối ưu giao diện responsive.                       |
| Backend        | Java Spring Boot hoặc Node.js/NestJS        | Xử lý nghiệp vụ và API.                            |
| Database       | PostgreSQL hoặc MySQL                       | Lưu dữ liệu quan hệ.                               |
| QR Library     | Thư viện tạo QR phía backend hoặc frontend  | Sinh mã QR theo nội dung và cấu hình.              |
| Payment        | VNPay, Momo, ZaloPay, PayPal hoặc Stripe    | Tích hợp thanh toán Pro.                           |
| Analytics      | User-Agent parser, IP geolocation tương đối | Phân tích thiết bị và vị trí ở mức không nhạy cảm. |

# **3. ĐỐI TƯỢNG SỬ DỤNG VÀ PHÂN QUYỀN**

## **3.1. Danh sách đối tượng sử dụng**

| **Đối tượng**   | **Mô tả**                           | **Mục tiêu sử dụng**                                                    |
|-----------------|-------------------------------------|-------------------------------------------------------------------------|
| Khách truy cập  | Người chưa đăng nhập vào hệ thống.  | Tìm hiểu website, tạo thử QR cơ bản, đăng ký tài khoản.                 |
| Người dùng Free | Người đã có tài khoản miễn phí.     | Tạo QR cơ bản, lưu QR giới hạn, tải PNG.                                |
| Người dùng Pro  | Người dùng đã nâng cấp gói trả phí. | Tạo Dynamic QR, tùy chỉnh nâng cao, xem thống kê, tải nhiều định dạng.  |
| Admin           | Người quản trị hệ thống.            | Quản lý người dùng, QR, gói, thanh toán, template và thống kê hệ thống. |
| Người quét QR   | Người dùng thiết bị quét mã QR.     | Truy cập nội dung QR và phát sinh dữ liệu lượt quét.                    |
| Cổng thanh toán | Tác nhân bên ngoài xử lý giao dịch. | Xác nhận thanh toán, trả kết quả giao dịch cho hệ thống.                |

## **3.2. Ma trận phân quyền tổng quát**

| **Chức năng**                  | **Khách**   | **Free**    | **Pro** | **Admin**     |
|--------------------------------|-------------|-------------|---------|---------------|
| Xem trang chủ/bảng giá         | Có          | Có          | Có      | Có            |
| Tạo QR cơ bản                  | Có giới hạn | Có          | Có      | Có            |
| Lưu QR vào tài khoản           | Không       | Có giới hạn | Có      | Có            |
| Tạo Dynamic QR                 | Không       | Không       | Có      | Có            |
| Chỉnh sửa Dynamic QR           | Không       | Không       | Có      | Có            |
| Thêm logo/template Pro         | Không       | Không       | Có      | Có            |
| Xem thống kê lượt quét         | Không       | Không       | Có      | Có            |
| Nâng cấp/gia hạn Pro           | Không       | Có          | Có      | Không áp dụng |
| Quản lý người dùng             | Không       | Không       | Không   | Có            |
| Quản lý gói/template/giao dịch | Không       | Không       | Không   | Có            |

# **4. YÊU CẦU CHỨC NĂNG**

Các yêu cầu chức năng được mã hóa theo dạng FR-XX để thuận tiện cho kiểm
thử và truy vết yêu cầu.

## **4.1. Nhóm FR-AUTH - Xác thực và tài khoản**

| **Mã yêu cầu** | **Tên yêu cầu**   | **Tác nhân**           | **Mô tả**                                                                                             | **Ưu tiên** |
|----------------|-------------------|------------------------|-------------------------------------------------------------------------------------------------------|-------------|
| FR-AUTH-01     | Đăng ký tài khoản | Khách truy cập         | Hệ thống cho phép đăng ký tài khoản bằng họ tên, email, số điện thoại, mật khẩu và xác nhận mật khẩu. | Must        |
| FR-AUTH-02     | Đăng nhập         | Khách/Người dùng/Admin | Hệ thống cho phép đăng nhập bằng email và mật khẩu hợp lệ.                                            | Must        |
| FR-AUTH-03     | Đăng xuất         | Người dùng/Admin       | Hệ thống cho phép người dùng đăng xuất và hủy phiên đăng nhập/token hiện tại.                         | Must        |
| FR-AUTH-04     | Cập nhật hồ sơ    | Free/Pro               | Người dùng có thể cập nhật họ tên, số điện thoại và ảnh đại diện nếu có.                              | Should      |
| FR-AUTH-05     | Đổi mật khẩu      | Free/Pro/Admin         | Người dùng có thể đổi mật khẩu khi cung cấp đúng mật khẩu cũ.                                         | Must        |
| FR-AUTH-06     | Khóa tài khoản    | Admin                  | Admin có thể khóa tài khoản vi phạm hoặc tài khoản cần tạm ngưng.                                     | Must        |

## **4.2. Nhóm FR-QR - Tạo mã QR**

| **Mã yêu cầu** | **Tên yêu cầu**           | **Tác nhân**   | **Mô tả**                                                              | **Ưu tiên** |
|----------------|---------------------------|----------------|------------------------------------------------------------------------|-------------|
| FR-QR-01       | Chọn loại QR              | Khách/Free/Pro | Hệ thống hiển thị danh sách loại QR được hỗ trợ và form tương ứng.     | Must        |
| FR-QR-02       | Tạo QR URL                | Khách/Free/Pro | Người dùng nhập URL hợp lệ để tạo QR dẫn đến website.                  | Must        |
| FR-QR-03       | Tạo QR Text               | Khách/Free/Pro | Người dùng nhập văn bản để tạo QR chứa nội dung text.                  | Must        |
| FR-QR-04       | Tạo QR WiFi               | Free/Pro       | Người dùng nhập SSID, mật khẩu và kiểu bảo mật để tạo QR kết nối WiFi. | Must        |
| FR-QR-05       | Tạo QR vCard              | Free/Pro       | Người dùng nhập thông tin liên hệ để tạo QR danh thiếp điện tử.        | Should      |
| FR-QR-06       | Tạo QR Email/SMS/Location | Free/Pro       | Người dùng tạo QR theo email, tin nhắn SMS hoặc vị trí bản đồ.         | Should      |
| FR-QR-07       | Tạo QR PDF/Menu/Social    | Pro            | Người dùng Pro có thể tạo các loại QR nâng cao.                        | Could       |

## **4.3. Nhóm FR-DESIGN - Tùy chỉnh thiết kế QR**

| **Mã yêu cầu** | **Tên yêu cầu**              | **Tác nhân** | **Mô tả**                                                                         | **Ưu tiên** |
|----------------|------------------------------|--------------|-----------------------------------------------------------------------------------|-------------|
| FR-DESIGN-01   | Tùy chỉnh màu sắc            | Free/Pro     | Người dùng có thể chọn màu QR và màu nền.                                         | Must        |
| FR-DESIGN-02   | Tùy chỉnh kích thước         | Free/Pro     | Người dùng có thể chọn kích thước QR trước khi tải xuống.                         | Should      |
| FR-DESIGN-03   | Tùy chỉnh kiểu mắt/kiểu chấm | Pro          | Người dùng Pro có thể chọn kiểu mắt QR và kiểu chấm nâng cao.                     | Should      |
| FR-DESIGN-04   | Thêm logo                    | Pro          | Người dùng Pro có thể tải logo lên và đặt vào giữa QR.                            | Should      |
| FR-DESIGN-05   | Chọn template                | Free/Pro     | Người dùng có thể chọn template miễn phí; template Pro chỉ mở cho người dùng Pro. | Should      |

## **4.4. Nhóm FR-MANAGE - Quản lý mã QR**

| **Mã yêu cầu** | **Tên yêu cầu**  | **Tác nhân** | **Mô tả**                                                          | **Ưu tiên** |
|----------------|------------------|--------------|--------------------------------------------------------------------|-------------|
| FR-MANAGE-01   | Lưu QR           | Free/Pro     | Người dùng đăng nhập có thể lưu mã QR vào tài khoản.               | Must        |
| FR-MANAGE-02   | Xem danh sách QR | Free/Pro     | Hệ thống hiển thị danh sách QR thuộc tài khoản đang đăng nhập.     | Must        |
| FR-MANAGE-03   | Tìm kiếm/lọc QR  | Free/Pro     | Người dùng có thể tìm QR theo tên, loại, trạng thái hoặc ngày tạo. | Should      |
| FR-MANAGE-04   | Xem chi tiết QR  | Free/Pro     | Người dùng xem thông tin chi tiết của QR thuộc tài khoản mình.     | Must        |
| FR-MANAGE-05   | Xóa QR           | Free/Pro     | Người dùng có thể xóa mềm QR sau khi xác nhận.                     | Must        |
| FR-MANAGE-06   | Quản lý thư mục  | Pro          | Người dùng Pro có thể tạo thư mục và nhóm QR theo thư mục.         | Could       |
| FR-MANAGE-07   | Sao chép QR      | Pro          | Người dùng Pro có thể tạo bản sao từ QR đã có.                     | Could       |

## **4.5. Nhóm FR-DYNAMIC - Dynamic QR và chuyển hướng**

| **Mã yêu cầu** | **Tên yêu cầu**       | **Tác nhân**  | **Mô tả**                                                                       | **Ưu tiên** |
|----------------|-----------------------|---------------|---------------------------------------------------------------------------------|-------------|
| FR-DYNAMIC-01  | Tạo Dynamic QR        | Pro           | Hệ thống tạo QR chứa đường dẫn trung gian dạng /q/{short_code}.                 | Must        |
| FR-DYNAMIC-02  | Chỉnh sửa URL đích    | Pro           | Chủ sở hữu QR có thể cập nhật URL đích của Dynamic QR.                          | Must        |
| FR-DYNAMIC-03  | Chuyển hướng khi quét | Người quét QR | Hệ thống nhận short_code, kiểm tra QR và chuyển hướng đến destination_url.      | Must        |
| FR-DYNAMIC-04  | Xử lý QR không hợp lệ | Người quét QR | Nếu QR bị xóa, khóa hoặc không tồn tại, hệ thống hiển thị trang lỗi thân thiện. | Must        |

## **4.6. Nhóm FR-ANALYTICS - Thống kê lượt quét**

| **Mã yêu cầu**  | **Tên yêu cầu**               | **Tác nhân** | **Mô tả**                                                                                                    | **Ưu tiên** |
|-----------------|-------------------------------|--------------|--------------------------------------------------------------------------------------------------------------|-------------|
| FR-ANALYTICS-01 | Ghi nhận lượt quét            | Hệ thống     | Mỗi lần Dynamic QR được quét, hệ thống ghi nhận thời gian, thiết bị, trình duyệt và vị trí tương đối nếu có. | Must        |
| FR-ANALYTICS-02 | Xem tổng lượt quét            | Pro          | Người dùng Pro xem tổng lượt quét của từng QR.                                                               | Must        |
| FR-ANALYTICS-03 | Biểu đồ theo thời gian        | Pro          | Hệ thống hiển thị lượt quét theo ngày, tuần hoặc tháng.                                                      | Should      |
| FR-ANALYTICS-04 | Thống kê thiết bị/trình duyệt | Pro          | Hệ thống phân nhóm lượt quét theo thiết bị, trình duyệt và hệ điều hành.                                     | Should      |
| FR-ANALYTICS-05 | Xuất báo cáo                  | Pro          | Người dùng Pro có thể xuất dữ liệu thống kê theo khoảng thời gian.                                           | Could       |

## **4.7. Nhóm FR-SUBPAY - Gói dịch vụ và thanh toán**

| **Mã yêu cầu** | **Tên yêu cầu**          | **Tác nhân**             | **Mô tả**                                                                              | **Ưu tiên** |
|----------------|--------------------------|--------------------------|----------------------------------------------------------------------------------------|-------------|
| FR-SUBPAY-01   | Xem bảng giá             | Khách/Free/Pro           | Hệ thống hiển thị thông tin gói Free và Pro.                                           | Must        |
| FR-SUBPAY-02   | Nâng cấp Pro             | Free                     | Người dùng Free có thể chọn gói Pro và tạo yêu cầu thanh toán.                         | Must        |
| FR-SUBPAY-03   | Xử lý kết quả thanh toán | Cổng thanh toán/Hệ thống | Hệ thống nhận callback/webhook và xác thực kết quả giao dịch.                          | Must        |
| FR-SUBPAY-04   | Kích hoạt Pro            | Hệ thống                 | Khi thanh toán thành công, hệ thống cập nhật subscription của người dùng thành ACTIVE. | Must        |
| FR-SUBPAY-05   | Gia hạn/hủy gói          | Pro                      | Người dùng Pro có thể gia hạn hoặc hủy gói theo chính sách hệ thống.                   | Should      |
| FR-SUBPAY-06   | Xem lịch sử thanh toán   | Free/Pro                 | Người dùng xem danh sách giao dịch của chính mình.                                     | Should      |

## **4.8. Nhóm FR-ADMIN - Quản trị hệ thống**

| **Mã yêu cầu** | **Tên yêu cầu**          | **Tác nhân** | **Mô tả**                                                                   | **Ưu tiên** |
|----------------|--------------------------|--------------|-----------------------------------------------------------------------------|-------------|
| FR-ADMIN-01    | Quản lý người dùng       | Admin        | Admin xem danh sách, tìm kiếm, xem chi tiết, khóa/mở khóa tài khoản.        | Must        |
| FR-ADMIN-02    | Quản lý QR toàn hệ thống | Admin        | Admin xem danh sách QR, lọc theo loại/trạng thái và vô hiệu hóa QR vi phạm. | Must        |
| FR-ADMIN-03    | Quản lý gói dịch vụ      | Admin        | Admin cấu hình gói Free/Pro, giá, giới hạn và quyền sử dụng.                | Should      |
| FR-ADMIN-04    | Quản lý giao dịch        | Admin        | Admin xem lịch sử giao dịch và xử lý giao dịch lỗi nếu cần.                 | Should      |
| FR-ADMIN-05    | Quản lý template QR      | Admin        | Admin tạo, cập nhật, ẩn hoặc xóa mềm template.                              | Should      |
| FR-ADMIN-06    | Xem dashboard hệ thống   | Admin        | Admin xem số người dùng, số QR, số lượt quét và doanh thu.                  | Should      |
| FR-ADMIN-07    | Xem log hệ thống         | Admin        | Admin xem nhật ký lỗi, cảnh báo và sự kiện bảo mật.                         | Should      |

# **5. ĐẶC TẢ USE CASE CHÍNH**

Các use case dưới đây mô tả những luồng nghiệp vụ quan trọng nhất của hệ
thống. Đây là cơ sở để thiết kế kiểm thử chức năng.

## **UC-01. Đăng ký tài khoản**

| **Thuộc tính** | **Nội dung**                      |
|----------------|-----------------------------------|
| Tác nhân chính | Khách truy cập                    |
| Tiền điều kiện | Khách truy cập chưa có tài khoản. |
| Hậu điều kiện  | Người dùng có tài khoản Free mới. |

**Luồng chính:**

1.  Người dùng mở form đăng ký.

2.  Người dùng nhập họ tên, email, số điện thoại, mật khẩu và xác nhận
    mật khẩu.

3.  Hệ thống kiểm tra dữ liệu hợp lệ và email chưa tồn tại.

4.  Hệ thống mã hóa mật khẩu và tạo tài khoản.

5.  Hệ thống gán vai trò USER và gói FREE mặc định.

Luồng ngoại lệ: Nếu email đã tồn tại hoặc dữ liệu không hợp lệ, hệ thống
hiển thị thông báo lỗi.

## **UC-02. Đăng nhập**

| **Thuộc tính** | **Nội dung**                                       |
|----------------|----------------------------------------------------|
| Tác nhân chính | Khách truy cập/Người dùng/Admin                    |
| Tiền điều kiện | Tài khoản tồn tại và chưa bị khóa.                 |
| Hậu điều kiện  | Người dùng vào được hệ thống theo quyền tương ứng. |

**Luồng chính:**

6.  Người dùng nhập email và mật khẩu.

7.  Hệ thống kiểm tra định dạng email.

8.  Hệ thống xác thực mật khẩu.

9.  Hệ thống kiểm tra trạng thái tài khoản.

10. Hệ thống tạo phiên đăng nhập hoặc token xác thực.

Luồng ngoại lệ: Nếu tài khoản bị khóa hoặc sai thông tin, hệ thống từ
chối đăng nhập.

## **UC-03. Tạo Static QR**

| **Thuộc tính** | **Nội dung**                                              |
|----------------|-----------------------------------------------------------|
| Tác nhân chính | Khách truy cập/Free/Pro                                   |
| Tiền điều kiện | Người dùng chọn loại QR và nhập dữ liệu hợp lệ.           |
| Hậu điều kiện  | Hệ thống sinh mã QR có thể xem trước, tải xuống hoặc lưu. |

**Luồng chính:**

11. Người dùng chọn loại QR.

12. Hệ thống hiển thị form nhập liệu phù hợp.

13. Người dùng nhập nội dung.

14. Hệ thống kiểm tra dữ liệu.

15. Người dùng tùy chỉnh thiết kế.

16. Hệ thống tạo QR và hiển thị bản xem trước.

17. Người dùng tải xuống hoặc lưu QR.

Luồng ngoại lệ: Nếu người dùng Free vượt giới hạn lưu QR, hệ thống yêu
cầu nâng cấp Pro.

## **UC-04. Tạo Dynamic QR**

| **Thuộc tính** | **Nội dung**                                     |
|----------------|--------------------------------------------------|
| Tác nhân chính | Người dùng Pro                                   |
| Tiền điều kiện | Người dùng đã đăng nhập và gói Pro còn hiệu lực. |
| Hậu điều kiện  | Dynamic QR được tạo với short_code và URL đích.  |

**Luồng chính:**

18. Người dùng chọn tạo Dynamic QR.

19. Người dùng nhập URL đích hoặc nội dung đích.

20. Hệ thống sinh short_code duy nhất.

21. Hệ thống tạo đường dẫn trung gian.

22. Hệ thống sinh mã QR chứa đường dẫn trung gian.

23. Người dùng lưu và tải QR.

Luồng ngoại lệ: Nếu gói Pro hết hạn, hệ thống không cho tạo Dynamic QR.

## **UC-05. Chỉnh sửa Dynamic QR**

| **Thuộc tính** | **Nội dung**                                       |
|----------------|----------------------------------------------------|
| Tác nhân chính | Người dùng Pro                                     |
| Tiền điều kiện | QR thuộc tài khoản người dùng và là Dynamic QR.    |
| Hậu điều kiện  | URL đích được cập nhật mà không cần tạo lại mã QR. |

**Luồng chính:**

24. Người dùng mở chi tiết Dynamic QR.

25. Người dùng chọn chỉnh sửa.

26. Người dùng nhập URL đích mới.

27. Hệ thống kiểm tra quyền sở hữu và dữ liệu hợp lệ.

28. Hệ thống cập nhật destination_url.

Luồng ngoại lệ: Static QR không được phép chỉnh sửa nội dung đã mã hóa.

## **UC-06. Quét Dynamic QR**

| **Thuộc tính** | **Nội dung**                                             |
|----------------|----------------------------------------------------------|
| Tác nhân chính | Người quét QR                                            |
| Tiền điều kiện | QR tồn tại, đang ACTIVE và có destination_url hợp lệ.    |
| Hậu điều kiện  | Người quét được chuyển hướng và lượt quét được ghi nhận. |

**Luồng chính:**

29. Người quét dùng thiết bị quét QR.

30. Trình duyệt truy cập đường dẫn /q/{short_code}.

31. Hệ thống tìm QR theo short_code.

32. Hệ thống ghi nhận lượt quét.

33. Hệ thống tăng scan_count.

34. Hệ thống chuyển hướng đến URL đích.

Luồng ngoại lệ: Nếu QR không tồn tại hoặc bị vô hiệu hóa, hệ thống hiển
thị trang lỗi.

## **UC-07. Nâng cấp gói Pro**

| **Thuộc tính** | **Nội dung**                                               |
|----------------|------------------------------------------------------------|
| Tác nhân chính | Người dùng Free                                            |
| Tiền điều kiện | Người dùng đã đăng nhập.                                   |
| Hậu điều kiện  | Tài khoản được nâng cấp Pro sau khi thanh toán thành công. |

**Luồng chính:**

35. Người dùng mở trang bảng giá.

36. Người dùng chọn gói Pro.

37. Hệ thống tạo giao dịch thanh toán.

38. Người dùng thanh toán qua cổng thanh toán.

39. Cổng thanh toán trả kết quả.

40. Hệ thống xác thực giao dịch.

41. Hệ thống cập nhật subscription thành ACTIVE.

Luồng ngoại lệ: Nếu thanh toán thất bại, tài khoản vẫn giữ gói Free.

## **UC-08. Admin quản lý người dùng**

| **Thuộc tính** | **Nội dung**                             |
|----------------|------------------------------------------|
| Tác nhân chính | Admin                                    |
| Tiền điều kiện | Admin đã đăng nhập.                      |
| Hậu điều kiện  | Admin xem và xử lý tài khoản người dùng. |

**Luồng chính:**

42. Admin mở trang quản lý người dùng.

43. Hệ thống hiển thị danh sách người dùng.

44. Admin tìm kiếm hoặc lọc tài khoản.

45. Admin xem chi tiết tài khoản.

46. Admin khóa hoặc mở khóa tài khoản nếu cần.

47. Hệ thống ghi log thao tác.

Luồng ngoại lệ: Admin không được xem mật khẩu dạng gốc của người dùng.

# **6. YÊU CẦU DỮ LIỆU VÀ THIẾT KẾ CƠ SỞ DỮ LIỆU**

## **6.1. Danh sách ENUM**

| **ENUM**               | **Giá trị**                                                     | **Mục đích**                           |
|------------------------|-----------------------------------------------------------------|----------------------------------------|
| user_status            | ACTIVE, LOCKED, DELETED                                         | Trạng thái tài khoản người dùng.       |
| role_name              | USER, ADMIN                                                     | Vai trò trong hệ thống.                |
| plan_name              | FREE, PRO                                                       | Tên gói dịch vụ.                       |
| plan_status            | ACTIVE, INACTIVE, DELETED                                       | Trạng thái gói dịch vụ.                |
| subscription_status    | ACTIVE, EXPIRED, CANCELLED, PENDING                             | Trạng thái đăng ký gói của người dùng. |
| payment_status         | PENDING, SUCCESS, FAILED, CANCELLED, REFUNDED                   | Trạng thái giao dịch thanh toán.       |
| payment_method         | VNPAY, MOMO, ZALOPAY, PAYPAL, STRIPE, BANK_TRANSFER             | Phương thức thanh toán.                |
| qr_type                | URL, TEXT, WIFI, VCARD, EMAIL, SMS, LOCATION, SOCIAL, PDF, MENU | Loại QR.                               |
| qr_status              | ACTIVE, DISABLED, DELETED                                       | Trạng thái mã QR.                      |
| template_status        | ACTIVE, HIDDEN, DELETED                                         | Trạng thái template QR.                |
| error_correction_level | L, M, Q, H                                                      | Mức sửa lỗi QR.                        |
| log_level              | INFO, WARNING, ERROR, SECURITY                                  | Mức độ nhật ký hệ thống.               |

## **6.2. Danh sách bảng dữ liệu**

| **STT** | **Tên bảng**  | **Mục đích**                             |
|---------|---------------|------------------------------------------|
| 1       | users         | Lưu thông tin tài khoản người dùng.      |
| 2       | roles         | Lưu danh sách vai trò.                   |
| 3       | user_roles    | Liên kết người dùng với vai trò.         |
| 4       | plans         | Lưu thông tin gói Free và Pro.           |
| 5       | subscriptions | Lưu gói dịch vụ người dùng đang sử dụng. |
| 6       | payments      | Lưu giao dịch thanh toán.                |
| 7       | folders       | Lưu thư mục quản lý QR.                  |
| 8       | qr_codes      | Lưu thông tin mã QR.                     |
| 9       | qr_designs    | Lưu cấu hình thiết kế QR.                |
| 10      | qr_templates  | Lưu template QR.                         |
| 11      | qr_scans      | Lưu lịch sử lượt quét QR.                |
| 12      | system_logs   | Lưu nhật ký hệ thống.                    |

## **6.3. Bảng users**

| **Tên trường** | **Kiểu dữ liệu** | **Ràng buộc**            | **Mô tả**                |
|----------------|------------------|--------------------------|--------------------------|
| id             | BIGINT           | PK, AUTO_INCREMENT       | Mã định danh người dùng. |
| full_name      | VARCHAR(100)     | NOT NULL                 | Họ và tên người dùng.    |
| email          | VARCHAR(150)     | NOT NULL, UNIQUE         | Email đăng nhập.         |
| password_hash  | VARCHAR(255)     | NOT NULL                 | Mật khẩu đã mã hóa.      |
| phone_number   | VARCHAR(20)      | NULL                     | Số điện thoại.           |
| avatar_url     | VARCHAR(255)     | NULL                     | Đường dẫn ảnh đại diện.  |
| status         | ENUM user_status | NOT NULL, DEFAULT ACTIVE | Trạng thái tài khoản.    |
| created_at     | DATETIME         | NOT NULL                 | Thời gian tạo.           |
| updated_at     | DATETIME         | NULL                     | Thời gian cập nhật.      |

## **6.4. Bảng roles**

| **Tên trường** | **Kiểu dữ liệu** | **Ràng buộc**      | **Mô tả**      |
|----------------|------------------|--------------------|----------------|
| id             | BIGINT           | PK, AUTO_INCREMENT | Mã vai trò.    |
| name           | ENUM role_name   | NOT NULL, UNIQUE   | Tên vai trò.   |
| description    | VARCHAR(255)     | NULL               | Mô tả vai trò. |

## **6.5. Bảng user_roles**

| **Tên trường** | **Kiểu dữ liệu** | **Ràng buộc**        | **Mô tả**      |
|----------------|------------------|----------------------|----------------|
| user_id        | BIGINT           | PK, FK -\> users(id) | Mã người dùng. |
| role_id        | BIGINT           | PK, FK -\> roles(id) | Mã vai trò.    |

## **6.6. Bảng plans**

| **Tên trường**       | **Kiểu dữ liệu** | **Ràng buộc**            | **Mô tả**              |
|----------------------|------------------|--------------------------|------------------------|
| id                   | BIGINT           | PK, AUTO_INCREMENT       | Mã gói dịch vụ.        |
| name                 | ENUM plan_name   | NOT NULL, UNIQUE         | Tên gói dịch vụ.       |
| price                | DECIMAL(12,2)    | NOT NULL, DEFAULT 0      | Giá gói.               |
| duration_days        | INT              | NULL                     | Thời hạn theo ngày.    |
| max_qr_codes         | INT              | NULL                     | Số QR tối đa được lưu. |
| allow_dynamic_qr     | BOOLEAN          | NOT NULL, DEFAULT FALSE  | Cho phép Dynamic QR.   |
| allow_logo           | BOOLEAN          | NOT NULL, DEFAULT FALSE  | Cho phép thêm logo.    |
| allow_analytics      | BOOLEAN          | NOT NULL, DEFAULT FALSE  | Cho phép thống kê.     |
| allow_svg_pdf_export | BOOLEAN          | NOT NULL, DEFAULT FALSE  | Cho phép tải SVG/PDF.  |
| description          | TEXT             | NULL                     | Mô tả gói.             |
| status               | ENUM plan_status | NOT NULL, DEFAULT ACTIVE | Trạng thái gói.        |
| created_at           | DATETIME         | NOT NULL                 | Thời gian tạo.         |
| updated_at           | DATETIME         | NULL                     | Thời gian cập nhật.    |

## **6.7. Bảng subscriptions**

| **Tên trường** | **Kiểu dữ liệu**         | **Ràng buộc**              | **Mô tả**           |
|----------------|--------------------------|----------------------------|---------------------|
| id             | BIGINT                   | PK, AUTO_INCREMENT         | Mã đăng ký gói.     |
| user_id        | BIGINT                   | NOT NULL, FK -\> users(id) | Mã người dùng.      |
| plan_id        | BIGINT                   | NOT NULL, FK -\> plans(id) | Mã gói.             |
| start_date     | DATETIME                 | NOT NULL                   | Ngày bắt đầu.       |
| end_date       | DATETIME                 | NULL                       | Ngày kết thúc.      |
| status         | ENUM subscription_status | NOT NULL, DEFAULT ACTIVE   | Trạng thái đăng ký. |
| auto_renew     | BOOLEAN                  | NOT NULL, DEFAULT FALSE    | Tự động gia hạn.    |
| created_at     | DATETIME                 | NOT NULL                   | Thời gian tạo.      |
| updated_at     | DATETIME                 | NULL                       | Thời gian cập nhật. |

## **6.8. Bảng payments**

| **Tên trường**   | **Kiểu dữ liệu**    | **Ràng buộc**                  | **Mô tả**                        |
|------------------|---------------------|--------------------------------|----------------------------------|
| id               | BIGINT              | PK, AUTO_INCREMENT             | Mã giao dịch.                    |
| user_id          | BIGINT              | NOT NULL, FK -\> users(id)     | Mã người dùng.                   |
| subscription_id  | BIGINT              | NULL, FK -\> subscriptions(id) | Mã đăng ký gói liên quan.        |
| amount           | DECIMAL(12,2)       | NOT NULL                       | Số tiền thanh toán.              |
| currency         | VARCHAR(10)         | NOT NULL, DEFAULT VND          | Đơn vị tiền tệ.                  |
| payment_method   | ENUM payment_method | NOT NULL                       | Phương thức thanh toán.          |
| transaction_code | VARCHAR(100)        | UNIQUE                         | Mã giao dịch từ cổng thanh toán. |
| status           | ENUM payment_status | NOT NULL, DEFAULT PENDING      | Trạng thái giao dịch.            |
| paid_at          | DATETIME            | NULL                           | Thời gian thanh toán thành công. |
| created_at       | DATETIME            | NOT NULL                       | Thời gian tạo giao dịch.         |

## **6.9. Bảng folders**

| **Tên trường** | **Kiểu dữ liệu** | **Ràng buộc**              | **Mô tả**           |
|----------------|------------------|----------------------------|---------------------|
| id             | BIGINT           | PK, AUTO_INCREMENT         | Mã thư mục.         |
| user_id        | BIGINT           | NOT NULL, FK -\> users(id) | Mã người dùng.      |
| name           | VARCHAR(100)     | NOT NULL                   | Tên thư mục.        |
| description    | VARCHAR(255)     | NULL                       | Mô tả thư mục.      |
| created_at     | DATETIME         | NOT NULL                   | Thời gian tạo.      |
| updated_at     | DATETIME         | NULL                       | Thời gian cập nhật. |

## **6.10. Bảng qr_codes**

| **Tên trường**  | **Kiểu dữ liệu** | **Ràng buộc**            | **Mô tả**                    |
|-----------------|------------------|--------------------------|------------------------------|
| id              | BIGINT           | PK, AUTO_INCREMENT       | Mã QR.                       |
| user_id         | BIGINT           | NULL, FK -\> users(id)   | Mã người tạo QR.             |
| folder_id       | BIGINT           | NULL, FK -\> folders(id) | Mã thư mục chứa QR.          |
| title           | VARCHAR(150)     | NOT NULL                 | Tên mã QR.                   |
| qr_type         | ENUM qr_type     | NOT NULL                 | Loại QR.                     |
| content         | TEXT             | NOT NULL                 | Nội dung gốc của QR.         |
| short_code      | VARCHAR(100)     | UNIQUE                   | Mã định danh cho Dynamic QR. |
| is_dynamic      | BOOLEAN          | NOT NULL, DEFAULT FALSE  | QR tĩnh hay QR động.         |
| destination_url | TEXT             | NULL                     | URL đích của Dynamic QR.     |
| qr_image_url    | VARCHAR(255)     | NULL                     | Đường dẫn ảnh QR đã tạo.     |
| scan_count      | BIGINT           | NOT NULL, DEFAULT 0      | Tổng số lượt quét.           |
| status          | ENUM qr_status   | NOT NULL, DEFAULT ACTIVE | Trạng thái QR.               |
| created_at      | DATETIME         | NOT NULL                 | Thời gian tạo.               |
| updated_at      | DATETIME         | NULL                     | Thời gian cập nhật.          |

## **6.11. Bảng qr_designs**

| **Tên trường**         | **Kiểu dữ liệu**            | **Ràng buộc**                         | **Mô tả**             |
|------------------------|-----------------------------|---------------------------------------|-----------------------|
| id                     | BIGINT                      | PK, AUTO_INCREMENT                    | Mã thiết kế QR.       |
| qr_code_id             | BIGINT                      | NOT NULL, UNIQUE, FK -\> qr_codes(id) | Mã QR tương ứng.      |
| template_id            | BIGINT                      | NULL, FK -\> qr_templates(id)         | Mẫu thiết kế sử dụng. |
| foreground_color       | VARCHAR(20)                 | NULL                                  | Màu chính của QR.     |
| background_color       | VARCHAR(20)                 | NULL                                  | Màu nền của QR.       |
| eye_style              | VARCHAR(50)                 | NULL                                  | Kiểu mắt QR.          |
| dot_style              | VARCHAR(50)                 | NULL                                  | Kiểu chấm QR.         |
| frame_style            | VARCHAR(50)                 | NULL                                  | Kiểu khung QR.        |
| logo_url               | VARCHAR(255)                | NULL                                  | Đường dẫn logo.       |
| size                   | INT                         | NULL                                  | Kích thước QR.        |
| error_correction_level | ENUM error_correction_level | NOT NULL, DEFAULT M                   | Mức sửa lỗi QR.       |
| created_at             | DATETIME                    | NOT NULL                              | Thời gian tạo.        |
| updated_at             | DATETIME                    | NULL                                  | Thời gian cập nhật.   |

## **6.12. Bảng qr_templates**

| **Tên trường**    | **Kiểu dữ liệu**     | **Ràng buộc**            | **Mô tả**                    |
|-------------------|----------------------|--------------------------|------------------------------|
| id                | BIGINT               | PK, AUTO_INCREMENT       | Mã template.                 |
| name              | VARCHAR(100)         | NOT NULL                 | Tên template.                |
| preview_image_url | VARCHAR(255)         | NULL                     | Ảnh xem trước template.      |
| config_json       | TEXT                 | NULL                     | Cấu hình template dạng JSON. |
| is_pro            | BOOLEAN              | NOT NULL, DEFAULT FALSE  | Template dành cho Pro.       |
| status            | ENUM template_status | NOT NULL, DEFAULT ACTIVE | Trạng thái template.         |
| created_at        | DATETIME             | NOT NULL                 | Thời gian tạo.               |
| updated_at        | DATETIME             | NULL                     | Thời gian cập nhật.          |

## **6.13. Bảng qr_scans**

| **Tên trường**   | **Kiểu dữ liệu** | **Ràng buộc**                 | **Mô tả**                       |
|------------------|------------------|-------------------------------|---------------------------------|
| id               | BIGINT           | PK, AUTO_INCREMENT            | Mã lượt quét.                   |
| qr_code_id       | BIGINT           | NOT NULL, FK -\> qr_codes(id) | Mã QR được quét.                |
| scanned_at       | DATETIME         | NOT NULL                      | Thời gian quét.                 |
| ip_address       | VARCHAR(100)     | NULL                          | IP đã xử lý hoặc ẩn danh.       |
| user_agent       | TEXT             | NULL                          | Thông tin trình duyệt/thiết bị. |
| device_type      | VARCHAR(50)      | NULL                          | Loại thiết bị.                  |
| browser          | VARCHAR(50)      | NULL                          | Trình duyệt.                    |
| operating_system | VARCHAR(50)      | NULL                          | Hệ điều hành.                   |
| country          | VARCHAR(100)     | NULL                          | Quốc gia.                       |
| city             | VARCHAR(100)     | NULL                          | Thành phố.                      |
| referer          | VARCHAR(255)     | NULL                          | Nguồn truy cập nếu có.          |

## **6.14. Bảng system_logs**

| **Tên trường** | **Kiểu dữ liệu** | **Ràng buộc**          | **Mô tả**                   |
|----------------|------------------|------------------------|-----------------------------|
| id             | BIGINT           | PK, AUTO_INCREMENT     | Mã log.                     |
| user_id        | BIGINT           | NULL, FK -\> users(id) | Người thực hiện thao tác.   |
| action         | VARCHAR(100)     | NOT NULL               | Hành động được ghi nhận.    |
| entity_type    | VARCHAR(100)     | NULL                   | Loại đối tượng bị tác động. |
| entity_id      | BIGINT           | NULL                   | ID đối tượng bị tác động.   |
| level          | ENUM log_level   | NOT NULL, DEFAULT INFO | Mức độ log.                 |
| message        | TEXT             | NULL                   | Nội dung log.               |
| ip_address     | VARCHAR(100)     | NULL                   | IP thực hiện thao tác.      |
| created_at     | DATETIME         | NOT NULL               | Thời gian ghi log.          |

## **6.15. Quan hệ dữ liệu chính**

| **Quan hệ**                   | **Mô tả**                                                    |
|-------------------------------|--------------------------------------------------------------|
| users 1 - n qr_codes          | Một người dùng có thể tạo nhiều mã QR.                       |
| users n - n roles             | Một người dùng có thể có nhiều vai trò thông qua user_roles. |
| users 1 - n subscriptions     | Một người dùng có thể có nhiều lịch sử đăng ký gói.          |
| plans 1 - n subscriptions     | Một gói dịch vụ có thể được nhiều người dùng đăng ký.        |
| subscriptions 1 - n payments  | Một đăng ký gói có thể phát sinh nhiều giao dịch.            |
| users 1 - n folders           | Một người dùng có thể tạo nhiều thư mục.                     |
| folders 1 - n qr_codes        | Một thư mục có thể chứa nhiều mã QR.                         |
| qr_codes 1 - 1 qr_designs     | Một mã QR có một cấu hình thiết kế.                          |
| qr_templates 1 - n qr_designs | Một template có thể được nhiều QR sử dụng.                   |
| qr_codes 1 - n qr_scans       | Một mã QR có thể có nhiều lượt quét.                         |

# **7. YÊU CẦU GIAO DIỆN**

## **7.1. Nguyên tắc thiết kế giao diện**

- Giao diện rõ ràng, dễ hiểu, phù hợp với người dùng không chuyên kỹ
  thuật.

- Quy trình tạo QR cần ngắn gọn theo các bước: chọn loại QR, nhập nội
  dung, tùy chỉnh, xem trước, tải/lưu.

- Hỗ trợ responsive trên desktop, tablet và mobile.

- Các chức năng Pro cần hiển thị rõ trạng thái bị khóa đối với người
  dùng Free.

- Thông báo lỗi cần cụ thể, giúp người dùng biết dữ liệu nào cần sửa.

## **7.2. Danh sách màn hình chính**

| **Màn hình**       | **Đối tượng sử dụng** | **Nội dung chính**                                                    |
|--------------------|-----------------------|-----------------------------------------------------------------------|
| Trang chủ          | Tất cả                | Giới thiệu hệ thống, nút tạo QR, loại QR phổ biến, bảng giá.          |
| Đăng ký/Đăng nhập  | Khách truy cập        | Form đăng ký, đăng nhập, kiểm tra dữ liệu.                            |
| Tạo QR             | Khách/Free/Pro        | Chọn loại QR, nhập dữ liệu, tùy chỉnh, xem trước.                     |
| Quản lý QR         | Free/Pro              | Danh sách QR đã tạo, tìm kiếm, lọc, thao tác tải/xóa/chỉnh sửa.       |
| Chi tiết QR        | Free/Pro              | Thông tin QR, preview, link, trạng thái, lượt quét.                   |
| Thống kê QR        | Pro                   | Biểu đồ lượt quét theo thời gian, thiết bị, trình duyệt, vị trí.      |
| Bảng giá           | Khách/Free/Pro        | So sánh gói Free và Pro, nút nâng cấp.                                |
| Thanh toán         | Free/Pro              | Thông tin gói, số tiền, phương thức thanh toán, trạng thái giao dịch. |
| Admin Dashboard    | Admin                 | Tổng quan người dùng, QR, lượt quét, doanh thu, log.                  |
| Quản lý người dùng | Admin                 | Danh sách, tìm kiếm, khóa/mở khóa, xem chi tiết tài khoản.            |

# **8. YÊU CẦU PHI CHỨC NĂNG**

| **Mã** | **Nhóm yêu cầu** | **Mô tả**                                                                     | **Tiêu chí đánh giá**                                                |
|--------|------------------|-------------------------------------------------------------------------------|----------------------------------------------------------------------|
| NFR-01 | Bảo mật          | Mật khẩu phải được mã hóa trước khi lưu; không lưu mật khẩu dạng văn bản gốc. | Kiểm tra database không có plain text password.                      |
| NFR-02 | Phân quyền       | Tài nguyên cá nhân chỉ được truy cập bởi chủ sở hữu hoặc Admin.               | User A không thể truy cập QR của User B.                             |
| NFR-03 | Hiệu năng        | Tạo QR và xem trước QR phải phản hồi nhanh.                                   | Thời gian phản hồi thông thường dưới 2 giây với dữ liệu hợp lệ.      |
| NFR-04 | Chuyển hướng     | Dynamic QR phải chuyển hướng ổn định.                                         | Tỷ lệ chuyển hướng thành công cao khi QR ACTIVE và URL hợp lệ.       |
| NFR-05 | Tính dễ dùng     | Giao diện tạo QR phải trực quan và ít bước.                                   | Người dùng mới có thể tạo QR cơ bản mà không cần hướng dẫn phức tạp. |
| NFR-06 | Tương thích      | Website hoạt động trên trình duyệt hiện đại và nhiều kích thước màn hình.     | Không vỡ layout trên desktop và mobile.                              |
| NFR-07 | Bảo trì          | Mã nguồn cần tổ chức theo module nghiệp vụ.                                   | Có module Auth, User, QR, Payment, Analytics, Admin.                 |
| NFR-08 | Mở rộng          | Dễ bổ sung loại QR, template, gói dịch vụ và cổng thanh toán mới.             | Không cần thay đổi lớn ở các module không liên quan.                 |
| NFR-09 | Riêng tư         | Dữ liệu lượt quét không nên lưu thông tin cá nhân nhạy cảm.                   | IP được xử lý hoặc ẩn danh nếu cần.                                  |
| NFR-10 | Ghi log          | Các lỗi và thao tác quan trọng cần được ghi log.                              | Log có action, level, message, thời gian và user nếu có.             |
| NFR-11 | Toàn vẹn dữ liệu | Thanh toán chỉ kích hoạt Pro khi giao dịch SUCCESS.                           | Không cập nhật Pro với giao dịch PENDING/FAILED.                     |
| NFR-12 | Khả dụng         | Hệ thống cần hạn chế gián đoạn ở chức năng tạo QR và redirect.                | Có xử lý lỗi thân thiện khi dịch vụ gặp sự cố.                       |

# **9. QUY TẮC NGHIỆP VỤ**

| **Mã** | **Quy tắc**                                                                    |
|--------|--------------------------------------------------------------------------------|
| BR-01  | Email người dùng là duy nhất trong hệ thống.                                   |
| BR-02  | Người dùng mới sau đăng ký được gán vai trò USER và gói FREE mặc định.         |
| BR-03  | Tài khoản LOCKED không được đăng nhập.                                         |
| BR-04  | Người dùng chỉ được quản lý QR thuộc tài khoản của mình.                       |
| BR-05  | Người dùng Free bị giới hạn số lượng QR được lưu theo plans.max_qr_codes.      |
| BR-06  | Dynamic QR chỉ dành cho người dùng Pro còn hiệu lực.                           |
| BR-07  | Chỉ Dynamic QR mới có short_code và destination_url để chuyển hướng.           |
| BR-08  | Static QR không cho phép chỉnh sửa nội dung đã mã hóa sau khi tạo.             |
| BR-09  | Thêm logo và tải SVG/PDF chỉ mở cho người dùng Pro nếu cấu hình gói cho phép.  |
| BR-10  | Chỉ QR trạng thái ACTIVE mới được chuyển hướng khi người quét truy cập.        |
| BR-11  | QR trạng thái DELETED không hiển thị trong danh sách người dùng.               |
| BR-12  | Lượt quét chỉ ghi nhận cho Dynamic QR hoặc QR có cơ chế redirect của hệ thống. |
| BR-13  | Thanh toán SUCCESS mới kích hoạt hoặc gia hạn gói Pro.                         |
| BR-14  | Khi gói Pro hết hạn, người dùng quay về giới hạn của gói Free.                 |
| BR-15  | Admin có quyền vô hiệu hóa QR hoặc khóa tài khoản vi phạm.                     |

# **10. RÀNG BUỘC, GIẢ ĐỊNH VÀ PHỤ THUỘC**

## **10.1. Ràng buộc**

- Hệ thống cần có kết nối Internet đối với các chức năng đăng nhập, lưu
  QR, thanh toán, thống kê và Dynamic QR.

- Dynamic QR phụ thuộc vào server redirect; nếu server ngừng hoạt động,
  QR động có thể không chuyển hướng được.

- Dữ liệu vị trí trong thống kê chỉ mang tính tương đối, không đảm bảo
  chính xác tuyệt đối.

- Một số chức năng như thanh toán và phân tích vị trí phụ thuộc vào dịch
  vụ bên thứ ba.

- QR có logo hoặc tùy chỉnh quá mạnh cần đảm bảo vẫn đủ khả năng quét.

## **10.2. Giả định**

- Người dùng có email hợp lệ để đăng ký và nhận thông báo.

- Người dùng Pro thanh toán qua các cổng thanh toán được hệ thống hỗ
  trợ.

- Admin có tài khoản riêng và được cấp quyền quản trị trước khi vận
  hành.

- Các thiết bị quét QR sử dụng camera hoặc ứng dụng quét QR phổ biến.

## **10.3. Phụ thuộc**

- Phụ thuộc vào thư viện sinh mã QR.

- Phụ thuộc vào cổng thanh toán để xác nhận giao dịch.

- Phụ thuộc vào dịch vụ gửi email nếu hệ thống có xác thực email hoặc
  thông báo hết hạn Pro.

- Phụ thuộc vào hạ tầng server, database và domain để vận hành Dynamic
  QR.

# **11. TIÊU CHÍ NGHIỆM THU**

| **STT** | **Tiêu chí**                               | **Kết quả mong đợi**                                                        |
|---------|--------------------------------------------|-----------------------------------------------------------------------------|
| 1       | Đăng ký tài khoản bằng email chưa tồn tại. | Tạo tài khoản thành công, gán USER và FREE.                                 |
| 2       | Đăng nhập bằng thông tin hợp lệ.           | Người dùng vào dashboard tương ứng.                                         |
| 3       | Đăng nhập bằng tài khoản bị khóa.          | Hệ thống từ chối đăng nhập.                                                 |
| 4       | Tạo QR URL với URL hợp lệ.                 | QR được tạo và có thể xem trước/tải xuống.                                  |
| 5       | Tạo QR với dữ liệu không hợp lệ.           | Hệ thống hiển thị lỗi rõ ràng.                                              |
| 6       | Người dùng Free vượt giới hạn lưu QR.      | Hệ thống chặn lưu và gợi ý nâng cấp Pro.                                    |
| 7       | Người dùng Free tạo Dynamic QR.            | Hệ thống từ chối và yêu cầu nâng cấp.                                       |
| 8       | Người dùng Pro tạo Dynamic QR.             | QR có short_code và đường dẫn trung gian.                                   |
| 9       | Quét Dynamic QR đang ACTIVE.               | Hệ thống ghi nhận lượt quét và chuyển hướng đúng URL.                       |
| 10      | Chỉnh sửa URL đích của Dynamic QR.         | QR cũ vẫn dùng được và chuyển hướng đến URL mới.                            |
| 11      | Thanh toán Pro thành công.                 | Subscription chuyển sang ACTIVE và mở khóa chức năng Pro.                   |
| 12      | Thanh toán thất bại.                       | Tài khoản giữ nguyên gói Free/Pro hiện tại.                                 |
| 13      | Admin khóa tài khoản người dùng.           | Người dùng không thể đăng nhập sau khi bị khóa.                             |
| 14      | Admin vô hiệu hóa QR.                      | QR không còn chuyển hướng và hiển thị trang lỗi thân thiện.                 |
| 15      | Người dùng Pro xem thống kê.               | Hiển thị tổng lượt quét và thống kê theo thời gian/thiết bị nếu có dữ liệu. |

# **12. PHỤ LỤC**

## **12.1. Quy ước mã yêu cầu**

| **Tiền tố**  | **Ý nghĩa**                                       |
|--------------|---------------------------------------------------|
| FR-AUTH      | Yêu cầu chức năng nhóm xác thực và tài khoản.     |
| FR-QR        | Yêu cầu chức năng nhóm tạo QR.                    |
| FR-DESIGN    | Yêu cầu chức năng nhóm thiết kế QR.               |
| FR-MANAGE    | Yêu cầu chức năng nhóm quản lý QR.                |
| FR-DYNAMIC   | Yêu cầu chức năng nhóm Dynamic QR.                |
| FR-ANALYTICS | Yêu cầu chức năng nhóm thống kê.                  |
| FR-SUBPAY    | Yêu cầu chức năng nhóm gói dịch vụ và thanh toán. |
| FR-ADMIN     | Yêu cầu chức năng nhóm quản trị.                  |
| NFR          | Yêu cầu phi chức năng.                            |
| BR           | Quy tắc nghiệp vụ.                                |
| UC           | Use Case.                                         |

## **12.2. Gợi ý nội dung QR theo từng loại**

| **Loại QR** | **Dữ liệu đầu vào bắt buộc**            | **Ghi chú**                                       |
|-------------|-----------------------------------------|---------------------------------------------------|
| URL         | Đường dẫn website                       | URL cần đúng định dạng.                           |
| TEXT        | Nội dung văn bản                        | Có thể giới hạn số ký tự.                         |
| WIFI        | Tên WiFi, loại bảo mật, mật khẩu nếu có | Dữ liệu tạo theo chuẩn kết nối WiFi.              |
| VCARD       | Họ tên hoặc số điện thoại               | Có thể bổ sung email, địa chỉ, công ty.           |
| EMAIL       | Email người nhận                        | Có thể bổ sung tiêu đề và nội dung email.         |
| SMS         | Số điện thoại                           | Có thể bổ sung nội dung tin nhắn.                 |
| LOCATION    | Tọa độ hoặc địa chỉ                     | Có thể mở qua ứng dụng bản đồ.                    |
| SOCIAL      | Đường dẫn mạng xã hội                   | Dành cho Facebook, Instagram, TikTok, LinkedIn... |
| PDF         | File PDF hoặc URL tài liệu              | Nên dành cho Pro.                                 |
| MENU        | URL menu hoặc dữ liệu menu              | Phù hợp nhà hàng/quán cà phê.                     |

## **12.3. Kết luận tài liệu**

Tài liệu SRS này mô tả đầy đủ yêu cầu của dự án QR Generator - QR Studio
ở mức phục vụ phân tích, thiết kế, triển khai và kiểm thử. Các nội dung
trong tài liệu có thể tiếp tục được mở rộng thành tài liệu thiết kế chi
tiết, đặc tả API, kế hoạch kiểm thử và tài liệu hướng dẫn sử dụng.

*Tài liệu được xây dựng dựa trên phạm vi và yêu cầu đã thảo luận cho dự
án website tạo QR-Code với hai gói Free và Pro.*
