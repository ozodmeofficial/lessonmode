#!/usr/bin/env python3
"""
LessonMode blocklist generator.

Bu skript bloklanadigan domenlar ro'yxatini (1000+) yaratadi.
Kategoriyalar: AI xizmatlari, ijtimoiy tarmoqlar, messenjerlar, video/stream,
o'yin va boshqa chalg'ituvchi saytlar. Natija: src/LessonMode/Resources/blocklist.txt

Qayta ishlatish: python3 tools/gen_blocklist.py
"""
import os

# --- Asosiy platformalar (aniq domenlar) ---
AI = [
    "openai.com", "chat.openai.com", "chatgpt.com", "api.openai.com", "oaistatic.com",
    "oaiusercontent.com", "sora.com", "sora.chatgpt.com",
    "claude.ai", "anthropic.com", "console.anthropic.com", "claude.com",
    "gemini.google.com", "bard.google.com", "aistudio.google.com", "makersuite.google.com",
    "deepmind.com", "deepmind.google",
    "copilot.microsoft.com", "bing.com/chat", "copilot.cloud.microsoft",
    "perplexity.ai", "pplx.ai",
    "poe.com", "quora.com",
    "character.ai", "beta.character.ai", "c.ai",
    "huggingface.co", "hf.co",
    "deepseek.com", "chat.deepseek.com",
    "mistral.ai", "chat.mistral.ai", "lechat.mistral.ai",
    "x.ai", "grok.com", "grok.x.ai",
    "meta.ai",
    "cohere.com", "cohere.ai",
    "jasper.ai", "copy.ai", "writesonic.com", "rytr.me", "sudowrite.com",
    "notion.so", "notion.ai",
    "phind.com", "you.com", "yeschat.ai", "chatpdf.com", "askyourpdf.com",
    "gamma.app", "tome.app", "beautiful.ai",
    "midjourney.com", "leonardo.ai", "stability.ai", "dreamstudio.ai",
    "runwayml.com", "runway.com", "pika.art", "kling.ai", "klingai.com",
    "elevenlabs.io", "play.ht", "murf.ai", "descript.com",
    "replit.com", "codeium.com", "tabnine.com", "cursor.com", "cursor.sh",
    "blackbox.ai", "useblackbox.io", "codegeex.cn",
    "janitorai.com", "crushon.ai", "chai.ml", "kindroid.ai", "replika.com", "replika.ai",
    "turnitin.com", "gptzero.me", "quillbot.com", "paraphrasing-tool.com",
    "wolframalpha.com", "mathway.com", "photomath.com", "symbolab.com",
    "caktus.ai", "studdy.ai", "studyx.ai", "smodin.io", "humata.ai",
    "scite.ai", "consensus.app", "elicit.com", "scholarai.io",
    "ideogram.ai", "civitai.com", "tensor.art", "getimg.ai", "playground.com",
    "chatbot.com", "botpress.com", "voiceflow.com",
    "together.ai", "fireworks.ai", "groq.com", "openrouter.ai",
    "suno.com", "suno.ai", "udio.com",
    "lumalabs.ai", "haiper.ai", "genmo.ai", "synthesia.io", "heygen.com",
]

SOCIAL = [
    "facebook.com", "fb.com", "fb.me", "fbcdn.net", "messenger.com", "m.me",
    "instagram.com", "cdninstagram.com", "instagr.am", "ig.me", "threads.net", "threads.com",
    "twitter.com", "x.com", "t.co", "twimg.com",
    "tiktok.com", "tiktokcdn.com", "tiktokv.com", "musical.ly", "byteoversea.com",
    "pinterest.com", "pinimg.com", "pin.it",
    "snapchat.com", "snap.com",
    "reddit.com", "redd.it", "redditstatic.com", "redditmedia.com",
    "tumblr.com", "vk.com", "vk.ru", "userapi.com", "ok.ru", "odnoklassniki.ru",
    "linkedin.com", "licdn.com", "lnkd.in",
    "quora.com", "medium.com",
    "9gag.com", "9cache.com", "imgur.com", "ifunny.co",
    "weibo.com", "weibo.cn", "sinaimg.cn", "qq.com", "qzone.qq.com", "weixin.qq.com",
    "mastodon.social", "bsky.app", "bskye.app", "bsky.social",
    "discord.com", "discord.gg", "discordapp.com", "discordapp.net", "discord.media",
    "tinder.com", "bumble.com", "badoo.com",
    "flickr.com", "deviantart.com", "myspace.com", "meetme.com", "tagged.com",
    "likee.video", "likee.com", "bigo.tv", "bigolive.tv",
    "clubhouse.com", "joinclubhouse.com",
]

MESSENGERS = [
    "telegram.org", "web.telegram.org", "t.me", "telegram.me", "telegra.ph", "tdesktop.com",
    "telegram.dog", "telesco.pe",
    "whatsapp.com", "whatsapp.net", "wa.me", "web.whatsapp.com",
    "viber.com", "signal.org", "wickr.com", "line.me", "naver.jp",
    "kakao.com", "kakaocdn.net", "wechat.com", "skype.com",
    "imo.im", "icq.com", "snapchat.com", "zalo.me", "kik.com",
    "slack.com", "mattermost.com", "rocket.chat",
]

VIDEO = [
    "youtube.com", "youtu.be", "youtube-nocookie.com", "ytimg.com", "yt.be",
    "googlevideo.com", "youtubekids.com", "youtubei.googleapis.com",
    "twitch.tv", "ttvnw.net", "jtvnw.net", "twitchcdn.net",
    "netflix.com", "nflxvideo.net", "nflximg.net", "nflxext.com",
    "hulu.com", "disneyplus.com", "disney-plus.net", "hbomax.com", "max.com",
    "primevideo.com", "amazonvideo.com",
    "dailymotion.com", "dai.ly", "vimeo.com", "vimeocdn.com",
    "rutube.ru", "ok.ru", "kinopoisk.ru", "ivi.ru", "okko.tv",
    "bilibili.com", "bilibili.tv", "iqiyi.com", "youku.com", "tudou.com",
    "rumble.com", "bitchute.com", "odysee.com", "kick.com",
    "crunchyroll.com", "funimation.com", "9anime.to", "gogoanime.tv", "zoro.to",
    "spotify.com", "scdn.co", "spotifycdn.com", "soundcloud.com", "sndcdn.com",
    "deezer.com", "pandora.com", "tidal.com",
]

GAMES = [
    "roblox.com", "rbxcdn.com", "roblox.cn",
    "minecraft.net", "mojang.com",
    "epicgames.com", "fortnite.com", "unrealengine.com",
    "steampowered.com", "steamcommunity.com", "steamstatic.com", "steamgames.com",
    "ea.com", "origin.com", "ubisoft.com", "uplay.com", "ubi.com",
    "battle.net", "blizzard.com", "riotgames.com", "leagueoflegends.com", "valorant.com",
    "miniclip.com", "poki.com", "crazygames.com", "y8.com", "kongregate.com",
    "addictinggames.com", "coolmathgames.com", "friv.com", "agame.com", "armorgames.com",
    "chess.com", "lichess.org", "newgrounds.com", "itch.io",
    "pubg.com", "pubgmobile.com", "garena.com", "freefiremobile.com",
    "supercell.com", "clashofclans.com", "clashroyale.com", "brawlstars.com",
    "genshin.hoyoverse.com", "hoyoverse.com", "mihoyo.com",
    "nintendo.com", "playstation.com", "xbox.com",
    "twitchmetrics.net", "op.gg", "tracker.gg",
]

MISC_DISTRACTION = [
    "netflix.com", "hotstar.com", "pornhub.com", "xvideos.com", "xnxx.com",
    "onlyfans.com", "chaturbate.com", "stripchat.com",
    "aliexpress.com", "amazon.com", "ebay.com", "wish.com", "temu.com",
    "shein.com", "etsy.com",
    "buzzfeed.com", "boredpanda.com", "distractify.com", "ranker.com",
    "pikabu.ru", "4chan.org", "8kun.top", "kiwifarms.net",
    "giphy.com", "tenor.com",
    "espn.com", "livescore.com", "flashscore.com", "365scores.com",
    "wattpad.com", "fanfiction.net", "archiveofourown.org",
]

CATEGORIES = {
    "AI": AI,
    "SOCIAL": SOCIAL,
    "MESSENGERS": MESSENGERS,
    "VIDEO": VIDEO,
    "GAMES": GAMES,
    "MISC": MISC_DISTRACTION,
}

# Hech qachon bloklanmaydigan domenlar (Telegram bot API ustoz kodini yuborishi uchun)
NEVER_BLOCK = {
    "api.telegram.org",  # bot API - MUHIM
}

# Mintaqaviy / qo'shimcha variantlar generatsiyasi (1000+ ga yetkazish uchun)
# Ko'plab platformalar ko'p ccTLD/subdomenlardan foydalanadi.
CCTLDS = [
    "com", "net", "org", "co", "io", "me", "tv", "ru", "uz", "kz", "ua", "tr",
    "in", "id", "br", "mx", "de", "fr", "es", "it", "pl", "nl", "jp", "kr",
    "cn", "hk", "tw", "au", "ca", "uk", "info", "biz", "app", "ai", "gg",
]

SUBDOMAIN_PREFIXES = ["www", "m", "web", "app", "api", "mobile", "login", "static", "cdn", "media", "www2"]

# Variant generatsiya qilinadigan "poydevor" nomlar (platforma ildizlari)
VARIANT_SEEDS = [
    "facebook", "instagram", "tiktok", "youtube", "twitter", "telegram", "whatsapp",
    "snapchat", "pinterest", "reddit", "twitch", "discord", "netflix", "roblox",
    "chatgpt", "openai", "spotify", "vk", "vimeo", "dailymotion", "likee", "bigo",
]


def build():
    seen = set()
    lines = []
    lines.append("# LessonMode blocklist")
    lines.append("# Avtomatik generatsiya qilingan (tools/gen_blocklist.py). Qo'lda tahrir qilmang.")
    lines.append("# Format: har qatorda bitta domen. '#' bilan boshlangan qatorlar - izoh.")
    lines.append("# O'qituvchi qo'shimcha domenlarni config.json dagi 'extraBlocked' ro'yxatiga qo'shishi mumkin.")
    lines.append("")

    def add(domain, bucket):
        d = domain.strip().lower().strip("/")
        # URL path bo'lsa (masalan bing.com/chat) - faqat domen qismi
        if "/" in d:
            d = d.split("/", 1)[0]
        if not d or d in seen or d in NEVER_BLOCK:
            return
        seen.add(d)
        bucket.append(d)

    for cat, items in CATEGORIES.items():
        bucket = []
        for it in items:
            add(it, bucket)
        if bucket:
            lines.append(f"# === {cat} ({len(bucket)}) ===")
            lines.extend(sorted(bucket))
            lines.append("")

    # Variantlar (ccTLD + subdomen) - ro'yxatni kengaytirish uchun
    extra = []
    for seed in VARIANT_SEEDS:
        for tld in CCTLDS:
            add(f"{seed}.{tld}", extra)
        for pref in SUBDOMAIN_PREFIXES:
            add(f"{pref}.{seed}.com", extra)
    if extra:
        lines.append(f"# === VARIANTS ({len(extra)}) ===")
        lines.extend(sorted(extra))
        lines.append("")

    total = len(seen)
    lines.insert(4, f"# Jami domenlar: {total}")
    return "\n".join(lines) + "\n", total


def main():
    content, total = build()
    out = os.path.join(os.path.dirname(__file__), "..", "assets", "blocklist.txt")
    out = os.path.abspath(out)
    os.makedirs(os.path.dirname(out), exist_ok=True)
    with open(out, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"Yozildi: {out}")
    print(f"Jami domenlar: {total}")


if __name__ == "__main__":
    main()
