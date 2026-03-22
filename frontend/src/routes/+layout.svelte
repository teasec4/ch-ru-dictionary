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

<nav class="w-full bg-white text-black px-4 py-3 flex justify-center items-center">
    
    <!-- Десктопное меню -->
    <div class="hidden md:flex gap-6">
    {#each links as link}
        <a href={link.href} class="hover:text-gray-300">{link.text}</a>
    {/each}
    </div>
      
    <div class="md:hidden flex justify-between w-screen">
        <a href="/" class="px-4 py-2 hover:bg-gray-100 rounded font-bold">
            СЛОВАРЬ
        </a>
        <!-- Бургер кнопка на мобильных -->
        <button 
        class="focus:outline-none"
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
  <div class="md:hidden bg-white text-black flex flex-col gap-2 px-4 py-2">
    {#each links as link}
      <a href={link.href} class="py-2 hover:bg-blue-400 rounded" onclick={() => menuOpen = false}>{link.text}</a>
    {/each}
  </div>
{/if}

<main class="flex-1 w-full">
    {#if !menuOpen}
        {@render children()}
    {/if}
</main>

<style lang="postcss">
  @reference "tailwindcss";
</style>