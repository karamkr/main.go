# ✅ استخدم صورة Alpine الخفيفة مع الحزم المطلوبة
FROM golang:1.21-alpine  

# ✅ تثبيت الأدوات المطلوبة لـ SQLite و GORM
RUN apk add --no-cache gcc musl-dev sqlite-dev  

# ✅ تفعيل CGO
ENV CGO_ENABLED=1

# ✅ تحديد مجلد العمل داخل الحاوية
WORKDIR /app  

# ✅ نسخ ملفات المشروع إلى الحاوية
COPY . .  

# ✅ تنزيل التبعيات وبناء التطبيق
RUN go mod tidy && go build -o main .  

# ✅ تشغيل التطبيق عند تشغيل الحاوية
CMD ["/app/main"]
