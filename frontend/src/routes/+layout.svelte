<script lang="ts">
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
	<nav class="sticky top-0 z-50 w-full bg-white border-b border-dict-4">
		<div class="max-w-6xl mx-auto px-4 sm:px-6">
			<div class="flex justify-between items-center h-14 sm:h-16">
				<!-- Logo -->
				<a href="/" class="flex items-center gap-2 sm:gap-3">
					<div class="w-9 h-9 sm:w-10 sm:h-10 rounded-xl bg-gradient-primary flex items-center justify-center">
						<span class="text-white font-bold">字</span>
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
	
	<!-- Overlay -->
	{#if menuOpen}
		<!-- Overlay -->
		<div
			class="fixed inset-0 bg-black/30 z-40 md:hidden"
			onclick={closeMenu}
		/>

		<!-- Drawer -->
		<div
			class="fixed top-0 right-0 h-full w-4/5 max-w-sm bg-white z-50 shadow-xl
			       transform transition-transform duration-300
		        {menuOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}"
		>
			<div class="px-4 py-4 space-y-2">
				{#each links as link}
					<a
						href={link.href}
						class="block px-4 py-3 rounded-lg text-dict-2 hover:text-primary hover:bg-dict-4 transition"
						onclick={closeMenu}
					>
						{link.text}
					</a>
				{/each}

				<div class="pt-4 mt-4 border-t border-dict-4 text-xs text-dict-2 px-4">
					Версия 1.0.0
				</div>
			</div>
		</div>
		{/if}

	
	<main class="flex-1 w-full bg-linear-to-b from-primary-500 via-secondary-20 to-dict-4">
		{@render children()}
	</main>

	<!-- FOOTER -->
	<footer class="bg-white border-t border-dict-4 mt-10">
		<div class="max-w-6xl mx-auto px-4 sm:px-6 py-8">
			<div class="grid gap-6 sm:grid-cols-2 md:grid-cols-3">
				<!-- Brand -->
				<div>
					<div class="flex items-center gap-3 mb-3">
						<div class="w-8 h-8 rounded-lg bg-gradient-primary flex items-center justify-center">
							<span class="text-white text-sm font-bold">字</span>
						</div>
						<div>
							<div class="font-semibold text-text-primary text-sm">
								Китайский словарь
							</div>
							<div class="text-xs text-dict-2">汉俄词典</div>
						</div>
					</div>
					<p class="text-xs sm:text-sm text-dict-2">
						Современный китайско-русский словарь с удобным поиском.
					</p>
				</div>

				<!-- Nav -->
				<div>
					<h3 class="font-medium text-dict-1 mb-3 text-sm">Навигация</h3>
					<div class="space-y-1">
						{#each links as link}
							<a href={link.href} class="block text-sm text-dict-2 hover:text-primary">
								{link.text}
							</a>
						{/each}
					</div>
				</div>

				<!-- Support -->
				<div>
					<h3 class="font-medium text-dict-1 mb-3 text-sm">Поддержка</h3>
					<div class="space-y-1">
						<a href="#" class="block text-sm text-dict-2 hover:text-primary">Помощь</a>
						<a href="#" class="block text-sm text-dict-2 hover:text-primary">Обратная связь</a>
						<a href="#" class="block text-sm text-dict-2 hover:text-primary">Конфиденциальность</a>
					</div>
				</div>
			</div>

			<div class="mt-6 pt-6 border-t border-dict-4 text-center text-xs sm:text-sm text-dict-2">
				© 2026 Китайский словарь
			</div>
		</div>
	</footer>
</div>

<style lang="postcss">
	@reference "tailwindcss";
</style>