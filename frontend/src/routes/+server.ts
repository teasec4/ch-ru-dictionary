import { redirect } from '@sveltejs/kit';

export const GET = ({ request, cookies }) => {
    const savedLang = cookies.get('lang');
    if (savedLang === 'ru' || savedLang === 'zh') {
        throw redirect(302, `/${savedLang}`);
    }

    const header = request.headers.get('accept-language') || '';

    if (header.includes('zh')) {
        throw redirect(302, '/zh');
    }

    throw redirect(302, '/ru');
};