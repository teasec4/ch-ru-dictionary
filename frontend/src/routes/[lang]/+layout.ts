import { Language } from '$lib/context/language.js';
import { redirect } from '@sveltejs/kit';

export const load = ({ params }) => {
    const lang = params.lang;

    if (lang !== "ru" && lang !== "zh") {
        throw redirect(302, "/ru");
    }
  
  let translations = Language[lang];

    return { lang, translations };
};