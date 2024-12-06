Build:
docker build -t telegram-gobot .

Run: 
docker run -d --name telegram-gobot -e TELEGRAM_BOT_TOKEN=your_token telegram-gobot
