package ca.codeglacier.makishima.security;

import ca.codeglacier.makishima.auth.DiscordLoginSuccessHandler;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    private final DiscordLoginSuccessHandler successHandler;

    public SecurityConfig(DiscordLoginSuccessHandler successHandler) {
        this.successHandler = successHandler; // automatically injects through @Component annotation
    }

    @Bean
    public SecurityFilterChain oauthFilterChain(HttpSecurity http) throws Exception {
        return http
                .authorizeHttpRequests(auth -> auth.anyRequest().permitAll())
                .oauth2Login(oauth -> oauth.successHandler(successHandler))
                .build();
    }

}
