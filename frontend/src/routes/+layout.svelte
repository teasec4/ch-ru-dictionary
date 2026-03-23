<script lang="ts">
	let { children } = $props();
	import "../app.css";
	
	let menuOpen = $state(false);
	
	const toggleMenu = () => {
		menuOpen = !menuOpen;
	};
	
	const links = [
		{ href: "/", text: "Словарь" },
		{ href: "/about", text: "Новости" },
		{ href: "/contact", text: "Контакты" },
	]
</script>

<svelte:head>
    <title>Китайский словарь</title>
    <meta
        name="description"
        content="Современный китайско-русский словарь с поиском по иероглифам, пиньину и переводам"
    />

</svelte:head>

<div class="min-h-screen flex flex-col bg-bg-primary">

    <nav class="sticky top-0 z-50 w-full bg-white shadow-lg border-b border-dict-4">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between items-center h-16">
                 <div class="flex items-center">
                     <a href="/" class="flex items-center space-x-3 group">
                         <div class="w-10 h-10 rounded-xl bg-gradient-primary flex items-center justify-center shadow-md group-hover:shadow-lg transition-shadow">
                             <span>
                                 字
                             </span>
                         </div>
                         <div class="text-sm text-dict-2">
                             汉俄词典
                         </div>
                     </a>
                 </div>
                 
                 <!-- Десктопное меню -->
                 <div class="hidden md:flex items-center space-x-1">
                     {#each links as link}
                         <a
                             href={link.href}
                             class="flex items-center space-x-2 px-4 py-2 rounded-lg text-dict-2 hover:text-primary hover:bg-dict-4 transition-all duration-200 group"
                             
                         >
                            
                             <span>{link.text}</span>
                         </a>
                     {/each}
                 </div>
                 
                 <!-- мобильное меню -->
                 <div class="flex items-center space-x-4">
                     <!-- Бургер кнопка для мобильных -->
                     <button
                         class="md:hidden flex items-center justify-center w-10 h-10 rounded-lg hover:bg-dict-4 text-dict-2 hover:text-dict-1 transition-colors"
                         onclick={toggleMenu}
                         aria-label={menuOpen ? "Закрыть меню" : "Открыть меню"}
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

                     
                 </div>
            </div>
        </div>
       
       
    </nav> 
    
    <!-- Мобильное меню -->
    {#if menuOpen}
        <div
            class="md:hidden border-t border-dict-4 bg-white shadow-xl animate-slideDown"
        >

            <div class="px-4 py-3 space-y-1">
                {#each links as link}
                    <a
                        href={link.href}
                        class="flex items-center space-x-3 px-4 py-3 rounded-lg text-dict-2 hover:text-primary hover:bg-dict-4 transition-colors group"
                        onclick={() => (menuOpen = false)}
                        
                    >
                        
                        <span class="font-medium">{link.text}</span>
                    </a>
                {/each}

                <!-- Дополнительные элементы для мобильной версии -->
                <div class="pt-4 mt-4 border-t border-dict-4">
                    <div class="px-4 py-2 text-sm text-dict-2">
                        Версия 1.0.0
                    </div>
                </div>
            </div>

      </div>
    {/if}
    
    <main class="flex-1 w-full">
        {#if !menuOpen}
            {@render children()}
        {/if}
    </main>
    
    <!-- Футер -->
    <footer class="bg-white border-t border-dict-4 mt-12">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div class="grid md:grid-cols-3 gap-8">
                <div>
                    <div class="flex items-center space-x-3 mb-4">
                        <div
                            class="w-8 h-8 rounded-lg bg-gradient-primary flex items-center justify-center"
                        >
                            <span
                                class="chinese-text text-white text-sm font-bold"
                                >字</span
                            >
                        </div>
                        <div>
                            <div class="font-bold text-text-primary">
                                Китайский словарь
                            </div>
                            <div class="text-xs text-dict-2">汉俄词典</div>
                        </div>
                    </div>
                    <p class="text-sm text-dict-2">
                        Современный китайско-русский словарь с расширенным
                        поиском и удобным интерфейсом.
                    </p>
                </div>

                <div>
                    <h3 class="font-semibold text-dict-1 mb-4">Навигация</h3>
                    <div class="space-y-2">
                        {#each links.slice(0, 3) as link}
                            <a
                                href={link.href}
                                class="block text-sm text-dict-2 hover:text-primary transition-colors"
                            >
                                {link.text}
                            </a>
                        {/each}
                    </div>
                </div>

                <div>
                    <h3 class="font-semibold text-dict-1 mb-4">Поддержка</h3>
                    <div class="space-y-2">
                        <a
                            href="#"
                            class="block text-sm text-dict-2 hover:text-primary transition-colors"
                        >
                            Помощь
                        </a>
                        <a
                            href="#"
                            class="block text-sm text-dict-2 hover:text-primary transition-colors"
                        >
                            Обратная связь
                        </a>
                        <a
                            href="#"
                            class="block text-sm text-dict-2 hover:text-primary transition-colors"
                        >
                            Конфиденциальность
                        </a>
                    </div>
                </div>
            </div>

            <div class="mt-8 pt-8 border-t border-dict-4 text-center">
                <p class="text-sm text-dict-2">
                    © 2026 Китайский словарь. Все права защищены.
                </p>
            </div>
        </div>
    </footer>

</div>

<style lang="postcss">
  @reference "tailwindcss";
  
	
</style>