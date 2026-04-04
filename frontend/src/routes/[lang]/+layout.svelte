<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import LanguageToggle from "$lib/components/LanguageToggle.svelte";
    let { children, data } = $props();
    import "../../app.css";

    let menuOpen = $state(false);

    const toggleMenu = () => {
        menuOpen = !menuOpen;
    };

    const lang = data.lang;
    const translations = data.translations;
    console.log(lang);

    const closeMenu = () => (menuOpen = false);

    const links = [
        { href: `/${lang}`, text: translations.home },
        { href: `/${lang}/about`, text: translations.about },
        { href: `/${lang}/contact`, text: translations.contact },
        { href: `/${lang}/social`, text: translations.socialTitle || "Мы в соцсетях" },
        { href: `/${lang}/privacy`, text: translations.privacyTitle || "Права" },
    ];

    const isHomePage = page.url.pathname === `/${lang}` || page.url.pathname === `/${lang}/`;
    const containerClass = isHomePage ? "flex flex-col h-[100dvh]" : "flex flex-col";
</script>

<svelte:head>
    <title>{translations.siteTitle}</title>
    <meta name="description" content={translations.metaDescription} />
</svelte:head>

<div class={containerClass} onclick={() => menuOpen && closeMenu()}>
    <!-- NAVBAR -->
    <nav class="sticky top-0 z-50 w-full bg-white h-12">
        <div class="max-w-6xl mx-auto px-4 sm:px-6">
            <div class="flex justify-between items-center h-14">
                <!-- Logo -->
                <a
                    href="/"
                    onclick={(e) => {
                        e.preventDefault();
                        goto("/");
                    }}
                    class="flex items-center gap-2 sm:gap-3"
                >
                    <div
                        class="w-9 h-9 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center"
                    >
                        <span class=" font-bold">字</span>
                    </div>
                    <div class="text-xs sm:text-sm text-dict-2 leading-tight">
                        {translations.logoText}
                    </div>
                </a>

                <!-- Desktop menu -->
                <div class="hidden md:flex items-center gap-1">
                    {#each links as link}
                        <a
                            href={link.href}
                            class="px-3 py-2 rounded-lg text-sm text-dict-2 hover:text-primary hover:bg-dict-4 transition"
                        >
                            {link.text}
                        </a>
                    {/each}

                    <!-- Language toggle for desktop -->
                    <div class="ml-2">
                        <LanguageToggle />
                    </div>
                </div>

                <!-- Right section for mobile -->
                <div class="flex items-center gap-2 md:hidden relative">
                    <!-- Language toggle for mobile -->
                    <LanguageToggle />

                    <!-- Mobile menu button -->
                    <button
                        class="w-10 h-10 flex items-center justify-center rounded-lg hover:bg-dict-4 transition"
                        onclick={(e) => { e.stopPropagation(); toggleMenu(); }}
                        aria-label={menuOpen
                            ? translations.closeMenu
                            : translations.openMenu}
                    >
                        <svg
                            class="w-6 h-6"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            viewBox="0 0 24 24"
                        >
                            {#if !menuOpen}
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M4 6h16M4 12h16M4 18h16"
                                />
                            {:else}
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M6 18L18 6M6 6l12 12"
                                />
                            {/if}
                        </svg>
                    </button>

                    <!-- Mobile dropdown menu -->
                    {#if menuOpen}
                        <div class="absolute right-0 top-full mt-1 w-48 bg-white rounded-xl border border-dict-4 shadow-lg py-1 z-50">
                            {#each links as link}
                                <a
                                    href={link.href}
                                    onclick={closeMenu}
                                    class="block px-4 py-2 text-sm text-dict-2 hover:text-primary hover:bg-dict-4/30"
                                >
                                    {link.text}
                                </a>
                            {/each}
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </nav>

    <main class="flex-1 flex flex-col items-center justify-center w-full px-4 {isHomePage ? '' : ''}">
        {@render children()}
    </main>
</div>

<style lang="postcss">
    @reference "tailwindcss";
</style>
