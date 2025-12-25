# 🧅 Thor Scraper

**Thor Scraper**, Tor ağındaki `.onion` uzantılı siteleri otomatize bir şekilde ziyaret eden, erişilebilirlik durumlarını kontrol eden ve kanıt olarak **ekran görüntüsü (screenshot)** ile **HTML kaynak kodunu** kaydeden Go tabanlı bir CLI aracıdır.

Siber Tehdit İstihbaratı (CTI) süreçlerinde manuel olarak yapılan onion servis takibini hızlandırmak ve anonimliği koruyarak veri toplamak amacıyla geliştirilmiştir.

### Kullanım

**Tor Browser Açık Olmalı**

Programı çalıştırırken taramak istediğiniz dosyanın yolunu parametre olarak verin:

go run main.go <dosya_adi>

Örnekler:
go run main.go targets.yaml
go run main.go siteler.txt

## 🚀 Özellikler

* **Tam Anonimlik:** Tüm trafik ve DNS sorguları SOCKS5 (Port 9150) üzerinden tünellenir.
* **IP & DNS Sızıntı Koruması:** ChromeDP, yerel DNS çözümlemeyi engelleyecek şekilde yapılandırılmıştır.
* **Görsel Kanıt:** Sitelerin o anki durumunun `.png` formatında fotoğrafını çeker.
* **Veri Yedekleme:** Sayfanın `HTML` kodunu indirip saklar.
* **Hata Toleransı:** Yanıt vermeyen siteler için timeout mekanizması vardır, tarama kesilmez.

## 🛠️ Kurulum

Projeyi bilgisayarınıza klonlayın ve gerekli bağımlılıkları indirin:


git clone https://github.com/ferhatoozdemir14-dot/thor-scraper.git

cd thor-scraper

go mod tidy

go run main.go <dosya_adi>

