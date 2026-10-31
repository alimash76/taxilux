export type Lang = 'nl' | 'en';

export const strings = {
  nl: {
    siteName: '{{COMPANY_NAME}}',
    tagline: 'Taxidienst in {{COMPANY_CITY}}',
    requestRide: 'Rit aanvragen',
    callUs: 'Bel ons',
    callUsLabel: 'Bel {{COMPANY_PHONE}} voor korte of spoedritten',
    sendRequest: 'Aanvraag per e-mail',
    sendRequestLabel: 'Stuur een ritaanvraag voor lange of geplande ritten',
    placeholder: 'Deze website wordt nog ingericht. De definitieve content volgt zo spoedig mogelijk.',
    languageSwitch: 'English',
    notFound: 'Pagina niet gevonden',
  },
  en: {
    siteName: '{{COMPANY_NAME}}',
    tagline: 'Taxi service in {{COMPANY_CITY}}',
    requestRide: 'Request a ride',
    callUs: 'Call us',
    callUsLabel: 'Call {{COMPANY_PHONE}} for short or urgent trips',
    sendRequest: 'Send a request by email',
    sendRequestLabel: 'Submit a ride request for long or planned trips',
    placeholder: 'This website is still under construction. Final content will follow shortly.',
    languageSwitch: 'Nederlands',
    notFound: 'Page not found',
  },
} as const;

export function getStrings(lang: string) {
  return (lang === 'en' ? strings.en : strings.nl);
}
