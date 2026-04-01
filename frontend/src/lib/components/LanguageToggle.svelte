<script lang="ts">
    import { ChevronDown } from "lucide-svelte";
    import { page } from "$app/state";

    let isOpen = $state(false);

    const languages = [
        {
            code: "ru",
            name: "Русский",
            flag: "🇷🇺",
            label: "RU",
        },
        {
            code: "zh",
            name: "中文",
            flag: "🇨🇳",
            label: "ZH",
        },
    ];

    // Текущий язык из параметров страницы
    const currentLang = page.params.lang || "ru";
    const currentLanguage =
        languages.find((lang) => lang.code === currentLang) || languages[0];

    function toggleDropdown() {
        isOpen = !isOpen;
    }

    function closeDropdown() {
        isOpen = false;
    }

    function switchLanguage(langCode: string) {
        if (langCode === currentLang) {
            closeDropdown();
            return;
        }

        const currentPath = window.location.pathname;
        let newPath = "";

        // Обработка корневой страницы (/)
        if (currentPath === "/") {
            newPath = `/${langCode}`;
        }
        // Обработка страниц с языком в URL
        else if (currentPath.startsWith(`/${currentLang}`)) {
            const pathWithoutLang = currentPath.substring(
                `/${currentLang}`.length,
            );
            newPath = `/${langCode}${pathWithoutLang}`;
        }
        // Если язык не указан в URL
        else {
            newPath = `/${langCode}${currentPath}`;
        }

        closeDropdown();
        // Просто меняем URL - страница полностью перезагрузится
        window.location.href = newPath;
    }

    // Закрытие dropdown при клике вне компонента
    function handleClickOutside(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".language-toggle")) {
            closeDropdown();
        }
    }

    $effect(() => {
        if (isOpen) {
            document.addEventListener("click", handleClickOutside);
            return () =>
                document.removeEventListener("click", handleClickOutside);
        }
    });
</script>

<div class="language-toggle relative">
    <!-- Кнопка переключателя -->
    <button
        class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-dict-4/30 transition-all duration-200 text-sm text-dict-2 hover:text-primary group"
        onclick={toggleDropdown}
        aria-label="Переключить язык"
        aria-expanded={isOpen}
    >
        <span
            class="text-base transition-transform duration-200 group-hover:scale-110"
        >
            {currentLanguage.flag}
        </span>
        <span
            class="font-medium hidden sm:inline text-dict-1 group-hover:text-primary transition-colors duration-200"
        >
            {currentLanguage.label}
        </span>
        <ChevronDown
            class="w-3 h-3 transition-all duration-200 {isOpen
                ? 'rotate-180'
                : ''} text-dict-3 group-hover:text-primary"
        />
    </button>

    <!-- Dropdown меню -->
    {#if isOpen}
        <div
            class="absolute right-0 top-full mt-1 w-44 bg-white rounded-xl border border-dict-4 shadow-lg py-1 z-50 animate-fade-in overflow-hidden"
        >
            {#each languages as lang}
                <button
                    class="w-full px-4 py-2.5 flex items-center gap-3 hover:bg-dict-4/30 transition text-left group/item {lang.code ===
                    currentLang
                        ? 'bg-dict-4/20'
                        : ''}"
                    onclick={() => switchLanguage(lang.code)}
                    aria-current={lang.code === currentLang ? "true" : "false"}
                >
                    <span class="text-lg">{lang.flag}</span>
                    <div class="flex-1 min-w-0">
                        <div
                            class="font-medium text-dict-1 truncate group-hover/item:text-primary transition"
                        >
                            {lang.name}
                        </div>
                        <div class="text-xs text-dict-3 mt-0.5 truncate">
                            {lang.code === "ru" ? "Русский" : "中文"}
                        </div>
                    </div>
                    {#if lang.code === currentLang}
                        <div
                            class="w-1.5 h-1.5 rounded-full bg-accent shrink-0"
                        ></div>
                    {/if}
                </button>
            {/each}
        </div>
    {/if}
</div>

<style lang="postcss">
    @reference "tailwindcss";

    @keyframes fade-in {
        from {
            opacity: 0;
            transform: translateY(-4px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .animate-fade-in {
        animation: fade-in 0.15s ease-out;
    }
</style>
