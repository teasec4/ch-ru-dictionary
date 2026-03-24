<script lang="ts">
    import { goto } from "$app/navigation";
	let { children } = $props();
	import "../app.css";
	
	let menuOpen = $state(false);
	
	const toggleMenu = () => {
		menuOpen = !menuOpen;
	};
	
	const closeMenu = () => (menuOpen = false);
	
	const links = [
		{ href: "/", text: "Словарь" },
		{ href: "/about", text: "Новости" },
		{ href: "/contact", text: "Контакты" },
	];
	
	
</script>

<svelte:head>
	<title>Китайский словарь</title>
	<meta
		name="description"
		content="Современный китайско-русский словарь с поиском по иероглифам, пиньину и переводам"
	/>
</svelte:head>

<div class="min-h-screen flex flex-col">
        <!-- NAVBAR -->
	<nav class="sticky top-0 z-50 w-full bg-white">
		<div class="max-w-6xl mx-auto px-4 sm:px-6">
			<div class="flex justify-between items-center h-14">
				<!-- Logo -->
				<a href="/" onclick={(e) => {
                e.preventDefault();
                goto("/");
  }} class="flex items-center gap-2 sm:gap-3">
					<div class="w-9 h-9 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center">
						<span class=" font-bold">字</span>
					</div>
					<div class="text-xs sm:text-sm text-dict-2 leading-tight">
						汉俄词典
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
				</div>

				<!-- Mobile button -->
				<button
					class="md:hidden w-10 h-10 flex items-center justify-center rounded-lg hover:bg-dict-4 transition"
					onclick={toggleMenu}
					aria-label={menuOpen ? "Закрыть меню" : "Открыть меню"}
				>
					<svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
						{#if !menuOpen}
							<path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
						{:else}
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						{/if}
					</svg>
				</button>
			</div>
		</div>

	</nav>
	
	
	<main class="flex flex-1 justify-center items-center flex-col">
		{@render children()}
	</main>
	

	
</div>

<style lang="postcss">
	@reference "tailwindcss";
</style>