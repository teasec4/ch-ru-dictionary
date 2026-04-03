<script lang="ts">
    let { data } = $props();
    const t = data.translations;
    let message = $state("");
    let sent = $state(false);
    let sending = $state(false);

    async function sendFeedback() {
        if (!message.trim()) return;
        sending = true;
        await new Promise(r => setTimeout(r, 1000));
        sent = true;
        sending = false;
    }
</script>

<svelte:head>
    <title>{t.contactTitle} - {t.siteTitle}</title>
    <meta name="description" content={t.contactText} />
</svelte:head>

<div class="w-full max-w-2xl mx-auto py-8 px-4">
    <h1 class="text-2xl sm:text-3xl font-bold text-text-primary mb-6">{t.contactTitle}</h1>
    
    <div class="space-y-6">
        <p class="text-dict-2">{t.contactText}</p>
        
        <div class="p-4 rounded-xl bg-dict-4/20 border border-dict-4">
            <div class="flex items-center gap-2 mb-2">
                <span class="text-dict-3 text-sm">{t.email}:</span>
                <a href="mailto:feedback@chinesedict.com" class="text-accent hover:underline">feedback@chinesedict.com</a>
            </div>
        </div>

        {#if !sent}
            <div class="space-y-3">
                <textarea
                    bind:value={message}
                    placeholder={t.feedbackPlaceholder}
                    class="w-full p-3 rounded-xl border border-dict-3 focus:border-accent focus:ring-2 focus:ring-accent/20 outline-none resize-none h-32"
                ></textarea>
                <button
                    onclick={sendFeedback}
                    disabled={!message.trim() || sending}
                    class="px-6 py-2.5 rounded-xl bg-primary text-white font-semibold hover:bg-dict-1 disabled:opacity-50 transition"
                >
                    {sending ? t.loading : t.feedback}
                </button>
            </div>
        {:else}
            <div class="p-4 rounded-xl bg-green-50 border border-green-200">
                <p class="text-green-700">{t.feedbackSent}</p>
            </div>
        {/if}
    </div>
</div>