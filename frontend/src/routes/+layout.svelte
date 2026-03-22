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
	<title>Dictionary</title>
</svelte:head>

<div class="app">
    <nav class="bg-text-primary sticky top-0 z-50 w-full px-4 flex justify-center items-center">
        <!-- Десктопное меню -->
        <div class="hidden md:flex text-text-tertiary gap-15 py-5">
        {#each links as link}
            <a href={link.href}>{link.text}</a>
        {/each}
        </div>
          
        <div class="md:hidden flex justify-between items-center w-screen h-15">
            <a href="/" class="text-text-tertiary  px-4 font-bold">
                СЛОВАРЬ
            </a>
            <!-- Бургер кнопка на мобильных -->
            <button 
                class="focus:outline-none text-text-tertiary"
                onclick={toggleMenu}
            >
                <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                    {#if !menuOpen}
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
                    {:else}
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    {/if}
                </svg>
            </button>
        </div>
    </nav> 
    
    <!-- Мобильное меню -->
    {#if menuOpen}
      <div class="md:hidden flex flex-col gap-2 px-4 py-2">
        {#each links as link}
          <a href={link.href} class="py-2" onclick={() => menuOpen = false}>{link.text}</a>
        {/each}
      </div>
    {/if}
    
    <main class="flex-1 w-full">
        {#if !menuOpen}
            {@render children()}
        {/if}
    </main>
</div>

<style lang="postcss">
  @reference "tailwindcss";
  
	
</style>