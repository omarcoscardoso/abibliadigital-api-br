#!/usr/bin/env python3
import json
import os
import re

def parse_pug_to_html(loc):
    def get(key, default=""):
        return loc.get(key, default or key)

    html = f"""<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>ABíbliaDigital | {get("title")}</title>
  <meta name="author" content="@abibliadigital">
  <meta name="description" content="{get("description")}">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  
  <script defer src="/theme-v2/js/jquery-1.12.4.min.js"></script>
  <script defer src="/theme-v2/js/jquery.json-viewer.js"></script>
  <link rel="stylesheet" href="/theme-v2/css/jquery.json-viewer.css">
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@400;700&family=Playfair+Display:wght@700&display=swap">
  <link rel="stylesheet" href="/theme-v2/css/normalize.css">
  <link rel="stylesheet" href="/theme-v2/css/style.css">
  
  <link rel="icon" type="image/png" sizes="32x32" href="/theme-v2/images/favicon/favicon-32x32.png">
  <link rel="icon" type="image/png" sizes="16x16" href="/theme-v2/images/favicon/favicon-16x16.png">
  <link rel="manifest" href="/theme-v2/images/favicon/manifest.json">
  <meta name="theme-color" content="#ffffff">
</head>
<body>
  <!-- Navigation -->
  <nav>
    <div class="container">
      <a href="{get("thisLanguage")}"><img src="/theme-v2/images/logo-symbol-a-biblia-digital.svg" alt="{get("title")}"></a>
      <button class="icon-menu menu-toggle"></button>
      <ul>
        <li><button class="icon-close menu-toggle"></button></li>
        <li><a class="text section-scroll" href="#about" title="{get("aboutTitle")}">{get("about")}</a></li>
        <li><a class="text section-scroll" href="#howToUse" title="{get("howToUseTitle")}">{get("howToUse")}</a></li>
        <li><a class="text" href="https://github.com/omarcoscardoso/abibliadigital-api-br" target="_blank" title="{get("documentationTitle")}">{get("documentation")}</a></li>
        <li><a class="text section-scroll" href="#donate" title="{get("donateTitle")}">{get("donate")}</a></li>
        <li><a class="icon-github" href="https://github.com/omarcoscardoso/abibliadigital-api-br" target="_blank" title="{get("githubTitle")}"></a></li>
        <li><a class="icon-twitter" href="" target="_blank" title="{get("contactTitle")}"></a></li>
        <li><a class="icon-whatsapp" href="https://chat.whatsapp.com/D6u9LuqgKYHJwJXCnneXpH" target="_blank" title="{get("joinTheCommunity")}"></a></li>
        <li><a href="{get("changeLanguageLink")}" title="{get("changeLanguageTitle")}"><img src="{get("changeLanguageImage")}" alt="{get("changeLanguageTitle")}"></a></li>
      </ul>
    </div>
  </nav>

  <!-- Header -->
  <header>
    <a href="{get("thisLanguage")}"><img src="/theme-v2/images/logo-a-biblia-digital.svg" alt="{get("title")}"></a>
    <div class="container">
      <div id="hero">
        <div class="row">
          <div class="col">
            <h1>{get("heroTitle")}</h1>
            <h3>{get("heroSubtitle")}</h3>
            <div class="bg"></div>
          </div>
          <div class="col">
            <div id="counter">{get("counterTitle")}
              <small>{get("counterSubtitle")}</small>
              <div class="bg"></div>
              <div class="bg"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ul id="social-links">
      <li><a class="icon-github" href="https://github.com/omarcoscardoso/abibliadigital-api-br" target="_blank" title="{get("githubTitle")}"></a></li>
      <li><a class="icon-twitter" href="" target="_blank" title="{get("contactTitle")}"></a></li>
      <li><a class="icon-whatsapp" href="https://chat.whatsapp.com/D6u9LuqgKYHJwJXCnneXpH" target="_blank" title="{get("joinTheCommunity")}"></a></li>
    </ul>
  </header>

  <!-- About Section -->
  <section id="about">
    <div class="container">
      <div class="row">
        <div class="col">
          <h2>{get("manifestTitle")}</h2>
          <h5>{get("manifestSubtitle")}</h5>
          <p>{get("manifestDescription1")}</p>
          <p>{get("manifestDescription2")}</p>
          <p>{get("manifestDescription3")}</p>
        </div>
        <div class="col">
          <h4>{get("content")}</h4>
          <p>{get("contentDescription")}</p>
          <h4>{get("doc")}</h4>
          <p>{get("docDescription")}</p>
          <h4>{get("availability")}</h4>
          <p>{get("availabilityDescription")}</p>
          <h4>{get("search")}</h4>
          <p>{get("searchDescription")}</p>
        </div>
      </div>
    </div>
  </section>

  <!-- Banner CTA -->
  <section>
    <div class="container">
      <div class="banner">
        <div class="row">
          <div class="col">
            <h2>{get("bannerTitle")}</h2>
            <p>{get("bannerDescription")}</p>
          </div>
          <div class="col">
            <a class="btn-primary" href="https://abibliadigital.api.br/docs" target="_blank" title="{get("documentationTitle")}">{get("bannerButton")}</a>
            <ul id="documentation-links">
              <li><a class="doc-icon icon-svg" href="https://insomnia.rest/run/?label=AB%C3%ADbliaDigital%20API&uri=https%3A%2F%2Fraw%2Egithubusercontent%2Ecom%2Fomarcoscardoso%2Fabibliadigital-api-br%2Fmain%2Fdocs%2Finsomnia.json" target="_blank" title="{get("bannerDocInsomniaTitle")}"><img src="/theme-v2/images/ic-insomnia.svg" alt="{get("bannerDocInsomniaTitle")}"></a></li>
              <li><a class="doc-icon icon-svg" href="https://editor.swagger.io/?url=https%3A%2F%2Fraw%2Egithubusercontent%2Ecom%2Fomarcoscardoso%2Fabibliadigital-api-br%2Fmain%2Fdocs%2Fopenapi.yaml" target="_blank" title="{get("bannerDocSwaggerTitle")}"><img src="/theme-v2/images/ic-swagger.svg" alt="{get("bannerDocSwaggerTitle")}"></a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- How To Use -->
  <section id="howToUse">
    <div class="container">
      <h2>{get("howToUse")}</h2>
      <p class="font-p mg-bt-30">
        {get("options")}: books, books/gn, verses/nvi/gn/1, verses/acf/gn/1/1
        <a href="https://github.com/omarcoscardoso/abibliadigital-api-br/blob/main/DOCUMENTATION.md" target="_blank" title="{get("documentation")}">{get("andMore")}</a>
      </p>
      <form id="form-1">
        <label id="basic-addon3" for="path">https://abibliadigital.api.br/api/</label>
        <input id="path" type="text" aria-describedby="basic-addon3" value="verses/nvi/sl/23" placeholder="verses/{{version}}/{{book}}/{{chapter}}/{{number}}">
        <button id="button-addon1" type="submit" title="{get("search")}">{get("search")}</button>
      </form>
      <pre id="output"></pre>
    </div>
  </section>

  <!-- Donate -->
  <section id="donate">
    <div class="container">
      <h2>{get("contributeTitle")}</h2>
      <p class="font-p mg-bt-30">{get("contributeDescription")}</p>
    </div>
  </section>

  <!-- Footer -->
  <footer>
    <div class="container">
      <div class="row">
        <div class="col"><p>{get("credits")}</p></div>
        # <div class="col"><a class="btn-small" href="https://stats.uptimerobot.com/4YxIqI4OBm/803888373" target="_blank" title="{get("statusPage")}">{get("statusPage")}</a></div>
      </div>
    </div>
  </footer>

  <script defer src="/theme-v2/js/main.js"></script>
  <script defer src="/assets/scripts/app.js"></script>
</body>
</html>
"""
    return html

def main():
    os.makedirs("public/en", exist_ok=True)
    os.makedirs("public/pt", exist_ok=True)

    with open("data/locales/pt.json", "r", encoding="utf-8") as f:
        pt = json.load(f)

    with open("data/locales/en.json", "r", encoding="utf-8") as f:
        en = json.load(f)

    with open("public/index.html", "w", encoding="utf-8") as f:
        f.write(parse_pug_to_html(pt))

    with open("public/pt/index.html", "w", encoding="utf-8") as f:
        f.write(parse_pug_to_html(pt))

    with open("public/en/index.html", "w", encoding="utf-8") as f:
        f.write(parse_pug_to_html(en))

    print("Landing page HTML generated at public/index.html, public/pt/index.html and public/en/index.html!")

if __name__ == "__main__":
    main()
