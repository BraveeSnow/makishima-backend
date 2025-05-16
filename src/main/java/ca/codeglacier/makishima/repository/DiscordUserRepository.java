package ca.codeglacier.makishima.repository;

import ca.codeglacier.makishima.entity.DiscordUser;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface DiscordUserRepository extends JpaRepository<DiscordUser, Long> {
}
