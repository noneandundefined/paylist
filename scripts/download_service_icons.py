import io
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT_DIR = ROOT / "http-server" / "media" / "subscriptions"

USER_AGENT = "Mozilla/5.0 (compatible; PaylistIconBot/1.0)"

MISSING = [
    "hulu", "discovery-plus", "dazn", "espn-plus", "fubo", "sling-tv", "now-tv",
    "rakuten-viki", "mubi", "shudder", "curiosity-stream", "nebula", "more-tv",
    "amediateka", "megogo", "kion", "viju", "tvplus", "showjet", "amazon-music",
    "qobuz", "idagio", "zvuk", "boosteroid", "shadow-pc", "wow", "ffxiv", "eso",
    "lost-ark", "path-of-exile", "genshin-impact", "honkai-star-rail", "riot",
    "minecraft", "vk-play", "lesta", "box", "pcloud", "mega", "backblaze",
    "idrive", "sync-com", "mailru-cloud", "ticktick", "things", "obsidian",
    "lastpass", "dashlane", "keeper", "superhuman", "jira", "confluence",
    "setapp", "webflow", "convertkit", "hubspot", "calendly", "loom", "docusign",
    "quickbooks", "freshbooks", "toggl", "rescue-time", "raindrop", "pocket",
    "fastmail", "hey-email", "cyberghost", "pia", "adguard", "nextdns", "vercel",
    "netlify", "gitlab", "whatsapp", "x", "reddit", "ok", "nike-training",
    "fitbit", "myfitnesspal", "betterhelp", "whoop", "garmin", "freeletics",
    "future-fitness", "down-dog", "seven-workouts", "wsj", "washington-post",
    "the-athletic", "economist", "bloomberg", "ft", "wired", "meduza", "forbes",
    "rbc", "vedomosti", "kommersant", "the-bell", "skillshare", "masterclass",
    "brilliant", "babbel", "busuu", "rosetta-stone", "chess-com", "lichess",
    "udacity", "datacamp", "codecademy", "pluralsight", "skyeng",
    "skyeng-teachers", "lingualeo", "puzzle-english", "khan-academy",
    "epic-games", "battle-net", "luna", "zapier", "make", "ifttt", "bitdefender",
    "kaspersky", "norton", "mcafee", "avast", "malwarebytes", "any-do",
    "fantastical", "bear", "ulysses", "ia-writer", "craft", "roam-research",
    "mem", "reflect", "capcut", "descript", "epidemic-sound", "artlist",
    "envato", "shutterstock", "getty-images", "unsplash", "framer",
]

DOMAINS = {
    "hulu": "hulu.com",
    "discovery-plus": "discoveryplus.com",
    "dazn": "dazn.com",
    "espn-plus": "plus.espn.com",
    "fubo": "fubo.tv",
    "sling-tv": "sling.com",
    "now-tv": "nowtv.com",
    "rakuten-viki": "viki.com",
    "mubi": "mubi.com",
    "shudder": "shudder.com",
    "curiosity-stream": "curiositystream.com",
    "nebula": "nebula.tv",
    "more-tv": "more.tv",
    "amediateka": "amediateka.ru",
    "megogo": "megogo.net",
    "kion": "kion.ru",
    "viju": "viju.ru",
    "tvplus": "tvplus.ru",
    "showjet": "showjet.ru",
    "amazon-music": "music.amazon.com",
    "qobuz": "qobuz.com",
    "idagio": "idagio.com",
    "zvuk": "zvuk.com",
    "boosteroid": "boosteroid.com",
    "shadow-pc": "shadow.tech",
    "wow": "worldofwarcraft.blizzard.com",
    "ffxiv": "finalfantasyxiv.com",
    "eso": "elderscrollsonline.com",
    "lost-ark": "lostark.com",
    "path-of-exile": "pathofexile.com",
    "genshin-impact": "genshin.hoyoverse.com",
    "honkai-star-rail": "hsr.hoyoverse.com",
    "riot": "riotgames.com",
    "minecraft": "minecraft.net",
    "vk-play": "vkplay.ru",
    "lesta": "lesta.ru",
    "box": "box.com",
    "pcloud": "pcloud.com",
    "mega": "mega.nz",
    "backblaze": "backblaze.com",
    "idrive": "idrive.com",
    "sync-com": "sync.com",
    "mailru-cloud": "cloud.mail.ru",
    "ticktick": "ticktick.com",
    "things": "culturedcode.com",
    "obsidian": "obsidian.md",
    "lastpass": "lastpass.com",
    "dashlane": "dashlane.com",
    "keeper": "keepersecurity.com",
    "superhuman": "superhuman.com",
    "jira": "atlassian.com",
    "confluence": "atlassian.com",
    "setapp": "setapp.com",
    "webflow": "webflow.com",
    "convertkit": "kit.com",
    "hubspot": "hubspot.com",
    "calendly": "calendly.com",
    "loom": "loom.com",
    "docusign": "docusign.com",
    "quickbooks": "quickbooks.intuit.com",
    "freshbooks": "freshbooks.com",
    "toggl": "toggl.com",
    "rescue-time": "rescuetime.com",
    "raindrop": "raindrop.io",
    "pocket": "getpocket.com",
    "fastmail": "fastmail.com",
    "hey-email": "hey.com",
    "cyberghost": "cyberghostvpn.com",
    "pia": "privateinternetaccess.com",
    "adguard": "adguard.com",
    "nextdns": "nextdns.io",
    "vercel": "vercel.com",
    "netlify": "netlify.com",
    "gitlab": "gitlab.com",
    "whatsapp": "whatsapp.com",
    "x": "x.com",
    "reddit": "reddit.com",
    "ok": "ok.ru",
    "nike-training": "nike.com",
    "fitbit": "fitbit.com",
    "myfitnesspal": "myfitnesspal.com",
    "betterhelp": "betterhelp.com",
    "whoop": "whoop.com",
    "garmin": "garmin.com",
    "freeletics": "freeletics.com",
    "future-fitness": "future.co",
    "down-dog": "downdogapp.com",
    "seven-workouts": "workout.loseit.com",
    "wsj": "wsj.com",
    "washington-post": "washingtonpost.com",
    "the-athletic": "theathletic.com",
    "economist": "economist.com",
    "bloomberg": "bloomberg.com",
    "ft": "ft.com",
    "wired": "wired.com",
    "meduza": "meduza.io",
    "forbes": "forbes.com",
    "rbc": "rbc.ru",
    "vedomosti": "vedomosti.ru",
    "kommersant": "kommersant.ru",
    "the-bell": "thebell.io",
    "skillshare": "skillshare.com",
    "masterclass": "masterclass.com",
    "brilliant": "brilliant.org",
    "babbel": "babbel.com",
    "busuu": "busuu.com",
    "rosetta-stone": "rosettastone.com",
    "chess-com": "chess.com",
    "lichess": "lichess.org",
    "udacity": "udacity.com",
    "datacamp": "datacamp.com",
    "codecademy": "codecademy.com",
    "pluralsight": "pluralsight.com",
    "skyeng": "skyeng.ru",
    "skyeng-teachers": "skysmart.ru",
    "lingualeo": "lingualeo.com",
    "puzzle-english": "puzzle-english.com",
    "khan-academy": "khanacademy.org",
    "epic-games": "epicgames.com",
    "battle-net": "battle.net",
    "luna": "amazon.com",
    "zapier": "zapier.com",
    "make": "make.com",
    "ifttt": "ifttt.com",
    "bitdefender": "bitdefender.com",
    "kaspersky": "kaspersky.com",
    "norton": "norton.com",
    "mcafee": "mcafee.com",
    "avast": "avast.com",
    "malwarebytes": "malwarebytes.com",
    "any-do": "any.do",
    "fantastical": "flexibits.com",
    "bear": "bear.app",
    "ulysses": "ulysses.app",
    "ia-writer": "ia.net",
    "craft": "craft.do",
    "roam-research": "roamresearch.com",
    "mem": "mem.ai",
    "reflect": "reflect.app",
    "capcut": "capcut.com",
    "descript": "descript.com",
    "epidemic-sound": "epidemicsound.com",
    "artlist": "artlist.io",
    "envato": "envato.com",
    "shutterstock": "shutterstock.com",
    "getty-images": "gettyimages.com",
    "unsplash": "unsplash.com",
    "framer": "framer.com",
}


def already_exists(slug: str) -> bool:
    return any((OUT_DIR / f"{slug}{ext}").exists() for ext in (".png", ".jpg", ".jpeg", ".webp"))


def fetch(url: str) -> bytes | None:
    request = urllib.request.Request(url, headers={"User-Agent": USER_AGENT})
    try:
        with urllib.request.urlopen(request, timeout=20) as response:
            if response.status != 200:
                return None
            data = response.read()
            if not data or len(data) < 80:
                return None
            return data
    except (urllib.error.URLError, TimeoutError, ValueError):
        return None


def sources(domain: str) -> list[str]:
    return [
        f"https://logo.clearbit.com/{domain}",
        f"https://www.google.com/s2/favicons?sz=128&domain={domain}",
        f"https://icons.duckduckgo.com/ip3/{domain}.ico",
    ]


def main() -> int:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    ok = 0
    failed: list[str] = []

    for slug in MISSING:
        if already_exists(slug):
            print(f"skip {slug}")
            ok += 1
            continue

        domain = DOMAINS.get(slug, slug.replace("-", "") + ".com")
        payload = None
        for url in sources(domain):
            payload = fetch(url)
            if payload:
                break
            time.sleep(0.15)

        if payload is None:
            print(f"fail {slug} ({domain})")
            failed.append(slug)
            continue

        header = payload[:12]
        ext = ".png"
        if header.startswith(b"\xff\xd8"):
            ext = ".jpg"
        elif header.startswith(b"RIFF") and b"WEBP" in payload[:16]:
            ext = ".webp"
        elif header[:4] == b"\x00\x00\x01\x00":
            ext = ".ico"

        path = OUT_DIR / f"{slug}{ext}"
        path.write_bytes(payload)
        print(f"ok   {slug}{ext} ({len(payload)} bytes)")
        ok += 1
        time.sleep(0.12)

    print(f"done {ok}/{len(MISSING)}, failed={len(failed)}")
    if failed:
        print("failed:", ", ".join(failed))
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
