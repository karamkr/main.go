# ✅ استخدم صورة Alpine الخفيفة
FROM golang:1.21-alpine  

# ✅ تحديد مجلد العمل داخل الحاوية
WORKDIR /app  

# ✅ نسخ ملفات المشروع إلى الحاوية
COPY . .  

# ✅ تثبيت مكتبات SQLite المطلوبة
RUN apk add --no-cache gcc musl-dev sqlite-dev  

# ✅ تفعيل CGO أثناء البناء
ENV CGO_ENABLED=1  

# ✅ تنزيل التبعيات وبناء التطبيق
RUN go mod tidy && go build -o main .  

# ✅ تشغيل التطبيق عند تشغيل الحاوية
CMD ["/app/main"]
